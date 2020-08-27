import Content, { DefaultContent } from '@/api/models/Content';
import RestResponse from '@/api/models/RestResponse';
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
import SessionContent from '@/components/session/SessionContent';
import SessionForm from '@/components/session/SessionForm';
import SessionIndicator from '@/components/session/SessionIndicator';
import { NotificationHelper } from '@/utils/Notification';
import { useRouter } from 'next/router';
import React, { useCallback, useEffect, useRef, useState } from 'react';

const apiClient = createApiClient();
const lpClient = createLongPollingClient();
const sessionApi = new SessionApi(apiClient);
const qrImageApi = new QrImageApi(apiClient);

const notificationHelper = new NotificationHelper();
notificationHelper.initialize();

export default function SessionPage() {
  const router = useRouter();
  const handleError = (error: Error) => {
    console.error(error);
    router.push('/');
  };

  const [connectionState, setConnectionState] = useState(
    ConnectionState.DISCONNECTED
  );
  const [content, setContent] = useState(DefaultContent());

  const handleContent = useCallback((content: Content) => {
    setContent(content);
    notificationHelper.notify(document.title, content.payload);
  }, []);

  const setupApi = async (sessionId: string) => {
    try {
      const {
        data: { content },
      } = await sessionApi.getSession(sessionId);
      setContent(content);
    } catch (err) {
      handleError(err);
    }
  };

  const wsRef = useRef<WebSocket>();
  const setupWebSocket = useCallback(
    (sessionId: string) =>
      new Promise((resolve, reject) => {
        if (wsRef.current) {
          setConnectionState(ConnectionState.DISCONNECTED);
          wsRef.current.close();
        }

        const ws = createWebSocketClient(sessionId);
        ws.addEventListener('open', () => {
          setConnectionState(ConnectionState.CONNECTED);
        });
        ws.addEventListener('message', (evt) => {
          handleContent(JSON.parse(evt.data));
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
      }),
    []
  );

  const esRef = useRef<EventSource>();
  const setupSSE = useCallback(
    (sessionId: string) =>
      new Promise((_, reject) => {
        if (esRef.current) {
          setConnectionState(ConnectionState.DISCONNECTED);
          esRef.current.close();
        }

        const es = createEventSourceClient(sessionId);
        es.addEventListener('ping', () => {
          setConnectionState(ConnectionState.CONNECTED);
        });
        es.addEventListener('content', (evt: MessageEvent) => {
          handleContent(JSON.parse(evt.data));
        });
        es.addEventListener('error', reject);
        setConnectionState(ConnectionState.CONNECTING);

        esRef.current = es;
      }),
    []
  );

  const doLongPolling = useCallback(async (sessionId: string) => {
    setConnectionState(ConnectionState.CONNECTING);
    const resp = await lpClient.get(`/sessions/${sessionId}`);
    setConnectionState(ConnectionState.CONNECTED);
    const { data: content } = resp.data as RestResponse<Content>;
    handleContent(content);
    return doLongPolling(sessionId);
  }, []);

  const query = router.query;
  const sessionId = query.id as string;

  const handleReload = () => {
    if (!sessionId) return;
    setupApi(sessionId);
    setupWebSocket(sessionId)
      .catch((err) => {
        console.error(err);
        return setupSSE(sessionId);
      })
      .catch((err) => {
        console.error(err);
        return doLongPolling(sessionId);
      })
      .catch((err) => console.error(err));
  };

  useEffect(handleReload, [sessionId]);

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
        />
      </Card>
      <SessionForm api={sessionApi} sessionId={sessionId} />
      <SessionContent content={content} />
    </div>
  );
}
