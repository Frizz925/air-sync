import { DefaultContent } from '@/api/models/Content';
import SessionApi from '@/api/SessionApi';
import { createApiClient, createWebSocketClient } from '@/clients';
import ConnectionState from '@/components/models/ConnectionState';
import SessionActions from '@/components/SessionActions';
import SessionContent from '@/components/SessionContent';
import SessionForm from '@/components/SessionForm';
import SessionIndicator from '@/components/SessionIndicator';
import { useRouter } from 'next/router';
import React, { useEffect, useRef, useState } from 'react';

const sessionApi = new SessionApi(createApiClient());

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
  const setupWebSocket = (sessionId: string) => {
    if (wsRef.current) wsRef.current.close();
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
  };

  const query = router.query;
  const sessionId = query.id as string;

  useEffect(() => {
    if (!sessionId) return;
    setupApi(sessionId);
    setupWebSocket(sessionId);
  }, [sessionId]);

  return (
    <div className='container container-main space-y-4'>
      <div className='card'>
        <div className='py-2 px-4'>
          <SessionIndicator
            sessionId={sessionId}
            connectionState={connectionState}
          />
        </div>
        <SessionActions api={sessionApi} sessionId={sessionId} />
      </div>
      <SessionForm api={sessionApi} sessionId={sessionId} />
      <SessionContent content={content} />
    </div>
  );
}
