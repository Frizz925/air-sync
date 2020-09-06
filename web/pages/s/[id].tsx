import { MessageDeleted, MessageInserted } from '@/api/models/Events';
import Message from '@/api/models/Message';
import RestResponse from '@/api/models/RestResponse';
import SessionEvent from '@/api/models/SessionEvent';
import QrImageApi from '@/api/QrImageApi';
import SessionApi from '@/api/SessionApi';
import {
  createApiClient,
  createEventSourceClient,
  createLongPollingClient,
  createWebSocketClient,
} from '@/clients';
import Card from '@/components/common/Card';
import Confirm from '@/components/common/Confirm';
import ConnectionState from '@/components/models/ConnectionState';
import SessionActions from '@/components/session/SessionActions';
import SessionIndicator from '@/components/session/SessionIndicator';
import alert from '@/utils/Alert';
import { handleErrorAlert } from '@/utils/Error';
import { NotificationHelper } from '@/utils/Notification';
import { getAttachmentUrl } from '@/utils/Url';
import Axios, { AxiosError, CancelTokenSource } from 'axios';
import dynamic from 'next/dynamic';
import { useRouter } from 'next/router';
import { resolve } from 'path';
import React, { useEffect, useRef, useState } from 'react';

const SessionForm = dynamic(() => import('@/components/session/SessionForm'), {
  ssr: false,
});

const SessionMessages = dynamic(
  () => import('@/components/session/SessionMessages'),
  { ssr: false }
);

const { CancelToken } = Axios;

const apiClient = createApiClient();
const lpClient = createLongPollingClient();
const sessionApi = new SessionApi(apiClient);
const qrImageApi = new QrImageApi(apiClient);

const webSocketEnabled = process.env.NEXT_PUBLIC_WEBSOCKET_ENABLED === 'true';
const eventStreamEnabled =
  process.env.NEXT_PUBLIC_EVENT_STREAM_ENABLED === 'true';

const notificationHelper = new NotificationHelper();
notificationHelper.initialize();

export default function SessionPage() {
  const router = useRouter();

  const [connectionState, setConnectionState] = useState(
    ConnectionState.DISCONNECTED
  );
  const [messages, setMessages] = useState<Message[]>([]);
  const [timestamp, setTimestamp] = useState<number>(new Date().getTime());
  const [messageDialog, setMessageDialog] = useState(false);

  const messageRef = useRef<Message>();
  const cancelRef = useRef<CancelTokenSource>();

  const runningRef = useRef(true);
  const handleError = (error: Error) => {
    console.error(error);
    handleErrorAlert(error);
    runningRef.current = false;
    router.push('/');
  };

  const messagesRef = useRef<Message[]>([]);
  const handleMessage = (message: Message) => {
    const attachmentUrl =
      message.attachment_id && getAttachmentUrl(message.attachment_id);
    notificationHelper.notify(document.title, message.body, attachmentUrl);
    const newMessages = [message, ...messagesRef.current];
    setMessages(newMessages);
    messagesRef.current = newMessages;
  };

  const handleDeletedSession = () => {
    runningRef.current = false;
    router.push('/');
  };

  const handleInsertedMessage = (event: MessageInserted) => {
    handleMessage(event.message);
  };

  const handleDeletedMessage = (event: MessageDeleted) => {
    const newMessages = [...messagesRef.current].filter(
      (message) => message.id !== event.message_id
    );
    setMessages(newMessages);
    messagesRef.current = newMessages;
  };

  const handleSessionEvent = ({ event, data }: SessionEvent<any>) => {
    switch (event) {
      case 'session.deleted':
        handleDeletedSession();
        break;
      case 'message.inserted':
        handleInsertedMessage(data as MessageInserted);
        break;
      case 'message.deleted':
        handleDeletedMessage(data as MessageDeleted);
        break;
    }
  };

  const setupApi = async (sessionId: string) => {
    try {
      const {
        data: { messages },
      } = await sessionApi.getSession(sessionId);
      messagesRef.current = messages;
      setMessages(messages);
    } catch (err) {
      handleError(err);
    }
  };

  const wsRef = useRef<WebSocket>();
  const setupWebSocket = (sessionId: string) =>
    new Promise((resolve, reject) => {
      if (!webSocketEnabled) {
        reject('WebSocket disabled');
        return;
      }

      if (!runningRef.current) {
        resolve();
        return;
      }

      if (wsRef.current) {
        setConnectionState(ConnectionState.DISCONNECTED);
        wsRef.current.close();
      }

      const ws = createWebSocketClient(sessionId);
      ws.addEventListener('open', () => {
        setConnectionState(ConnectionState.CONNECTED);
      });
      ws.addEventListener('message', (e) => {
        handleSessionEvent(JSON.parse(e.data));
      });
      ws.addEventListener('error', (err) => {
        reject(err);
      });
      ws.addEventListener('close', () => {
        setConnectionState(ConnectionState.DISCONNECTED);
        resolve();
      });
      setConnectionState(ConnectionState.CONNECTING);

      wsRef.current = ws;
    });

  const esRef = useRef<EventSource>();
  const setupEventStream = (sessionId: string) =>
    new Promise((_, reject) => {
      if (!eventStreamEnabled) {
        reject('Event stream disabled');
        return;
      }

      if (!runningRef.current) {
        resolve();
        return;
      }

      if (esRef.current) {
        setConnectionState(ConnectionState.DISCONNECTED);
        esRef.current.close();
      }

      const es = createEventSourceClient(sessionId);
      es.addEventListener('heartbeat', () => {
        setConnectionState(ConnectionState.CONNECTED);
      });
      es.addEventListener('message', (e: MessageEvent) => {
        handleSessionEvent(JSON.parse(e.data));
      });
      es.addEventListener('close', () => {
        setConnectionState(ConnectionState.DISCONNECTED);
        es.close();
      });
      es.addEventListener('error', reject);
      setConnectionState(ConnectionState.CONNECTING);

      esRef.current = es;
    });

  const doLongPolling = async (sessionId: string) => {
    if (!runningRef.current) return;
    if (cancelRef.current) cancelRef.current.cancel('Long-polling restarted');

    let hasError = false;
    const start = new Date();
    const source = CancelToken.source();
    cancelRef.current = source;
    try {
      setConnectionState(ConnectionState.CONNECTED);
      const resp = await lpClient.get(`/sessions/${sessionId}`, {
        cancelToken: source.token,
      });
      if (resp.status === 200) {
        const event = (resp.data as RestResponse<SessionEvent<any>>).data;
        handleSessionEvent(event);
      }
    } catch (err) {
      if (Axios.isCancel(err)) {
        return Promise.resolve();
      }
      setConnectionState(ConnectionState.DISCONNECTED);
      console.error(err);
      handleErrorAlert(err);
      if (err.isAxiosError) {
        const resp = (err as AxiosError).response;
        if (resp && resp.status === 404) {
          return;
        }
      }
      hasError = true;
    }
    const end = new Date();

    return new Promise((resolve, reject) => {
      const diff = end.getTime() - start.getTime();
      // if last long-polling request didn't catch any error
      // or if it did catch error after 15 seconds
      // then execute the next request immediately
      if (!hasError || diff >= 15000) {
        doLongPolling(sessionId).then(resolve, reject);
        return;
      }
      // else wait for 3 seconds before executing another request
      setTimeout(() => {
        doLongPolling(sessionId).then(resolve, reject);
      }, 3000);
    });
  };

  const query = router.query;
  const sessionId = query.id as string;

  const doReload = () => {
    if (!sessionId) return;
    setConnectionState(ConnectionState.CONNECTING);
    setupApi(sessionId)
      .then(() => setupWebSocket(sessionId))
      .catch((err) => {
        console.error(err);
        return setupEventStream(sessionId);
      })
      .catch((err) => {
        console.error(err);
        return doLongPolling(sessionId);
      })
      .catch((err) => {
        console.error(err);
        handleErrorAlert(err);
      });
  };

  const handleReload = () => {
    alert('Reloading session...');
    doReload();
  };

  const handleDelete = () => {
    runningRef.current = false;
    router.push('/');
  };

  const deleteMessagePrompt = (message: Message) => {
    messageRef.current = message;
    setMessageDialog(true);
  };

  const handleMessageDelete = async () => {
    if (!messageRef.current) return;
    try {
      await sessionApi.deleteMessage(sessionId, messageRef.current.id);
      setMessageDialog(false);
    } catch (err) {
      console.error(err);
      handleErrorAlert(err);
    }
  };

  useEffect(() => {
    doReload();
    // Update timestamps every 30 seconds
    const interval = setInterval(() => {
      setTimestamp(new Date().getTime());
    }, 30000);
    return () => {
      clearInterval(interval);
    };
  }, [sessionId]);

  return (
    <div className='container-main space-y-4'>
      <Card>
        <div className='py-2 px-4'>
          <SessionIndicator
            sessionId={sessionId}
            connectionState={connectionState}
          />
        </div>
        <SessionActions
          sessionApi={sessionApi}
          sessionId={sessionId}
          qrImageApi={qrImageApi}
          onReload={handleReload}
          onDelete={handleDelete}
        />
      </Card>
      <SessionForm api={sessionApi} sessionId={sessionId} />
      <SessionMessages
        messages={messages}
        timestamp={timestamp}
        onDelete={deleteMessagePrompt}
      />
      <Confirm
        shown={messageDialog}
        message='Are you sure you want to delete this message?'
        confirmLabel='Delete'
        confirmColor='red'
        onConfirm={handleMessageDelete}
        onClose={() => setMessageDialog(false)}
      />
    </div>
  );
}
