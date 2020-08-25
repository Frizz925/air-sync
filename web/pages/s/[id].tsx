import { DefaultContent } from '@/api/models/Content';
import QrImageApi from '@/api/QrImageApi';
import SessionApi from '@/api/SessionApi';
import { createApiClient, createWebSocketClient } from '@/clients';
import Card from '@/components/common/Card';
import ConnectionState from '@/components/models/ConnectionState';
import SessionActions from '@/components/session/SessionActions';
import SessionContent from '@/components/session/SessionContent';
import SessionForm from '@/components/session/SessionForm';
import SessionIndicator from '@/components/session/SessionIndicator';
import { useRouter } from 'next/router';
import React, { useCallback, useEffect, useRef, useState } from 'react';

const apiClient = createApiClient();
const sessionApi = new SessionApi(apiClient);
const qrImageApi = new QrImageApi(apiClient);

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
  const setupWebSocket = useCallback((sessionId: string) => {
    if (wsRef.current) {
      setConnectionState(ConnectionState.DISCONNECTED);
      wsRef.current.close();
    }

    const ws = createWebSocketClient(sessionId);
    ws.addEventListener('open', () => {
      setConnectionState(ConnectionState.CONNECTED);
    });
    ws.addEventListener('message', (evt) => {
      setContent(JSON.parse(evt.data));
    });
    ws.addEventListener('error', (err) => {
      console.error(err);
    });
    ws.addEventListener('close', () => {
      setConnectionState(ConnectionState.DISCONNECTED);
    });
    setConnectionState(ConnectionState.CONNECTING);

    wsRef.current = ws;
  }, []);

  const query = router.query;
  const sessionId = query.id as string;

  const handleReload = () => {
    if (!sessionId) return;
    setupApi(sessionId);
    setupWebSocket(sessionId);
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
