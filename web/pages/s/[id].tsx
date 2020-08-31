import Message from '@/api/models/Message';
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
import ConnectionState from '@/components/models/ConnectionState';
import SessionActions from '@/components/session/SessionActions';
import SessionForm from '@/components/session/SessionForm';
import SessionIndicator from '@/components/session/SessionIndicator';
import SessionMessage from '@/components/session/SessionMessage';
import { NotificationHelper } from '@/utils/Notification';
import map from 'lodash/map';
import { useRouter } from 'next/router';
import { resolve } from 'path';
import React, { useEffect, useRef, useState } from 'react';

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
  const handleError = (error: Error) => {
    console.error(error);
    router.push('/');
  };

  const [running, setRunning] = useState(true);
  const [connectionState, setConnectionState] = useState(
    ConnectionState.DISCONNECTED
  );
  const [messages, setMessages] = useState<Message[]>([]);
  const [timestamp, setTimestamp] = useState<number>(new Date().getTime());

  const messagesRef = useRef<Message[]>([]);
  const handleMessage = (message: Message) => {
    notificationHelper.notify(document.title, message.content);
    const newMessages = [message, ...messagesRef.current];
    setMessages(newMessages);
    messagesRef.current = newMessages;
  };

  const handleDeletedMessage = (messageId: string) => {
    const newMessages = [...messagesRef.current].filter(
      (message) => message.id !== messageId
    );
    setMessages(newMessages);
    messagesRef.current = newMessages;
  };

  const handleSessionEvent = ({ event: name, data }: SessionEvent<any>) => {
    switch (name) {
      case 'message/insert':
        handleMessage(data as Message);
        break;
      case 'message/delete':
        handleDeletedMessage(data as string);
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

      if (!running) {
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

      if (!running) {
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
    if (!running) return;
    let hasError = false;
    const start = new Date();
    try {
      setConnectionState(ConnectionState.CONNECTED);
      const resp = await lpClient.get(`/sessions/${sessionId}`);
      if (resp.status === 200) {
        const event = resp.data as SessionEvent<unknown>;
        handleSessionEvent(event);
      }
    } catch (err) {
      setConnectionState(ConnectionState.DISCONNECTED);
      console.error(err);
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

  const handleReload = () => {
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
      .catch((err) => console.error(err));
  };

  const handleDelete = () => setRunning(false);

  useEffect(() => {
    handleReload();
    // Update timestamps every 30 seconds
    const interval = setInterval(() => {
      setTimestamp(new Date().getTime());
    }, 30000);
    return () => {
      clearInterval(interval);
    };
  }, [sessionId]);

  const messageComponents = map(messages, (message) => (
    <SessionMessage
      key={message.id}
      api={sessionApi}
      sessionId={sessionId}
      message={message}
      timestamp={timestamp}
    />
  ));

  return (
    <div className='container container-main space-y-4'>
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
      {messageComponents}
    </div>
  );
}
