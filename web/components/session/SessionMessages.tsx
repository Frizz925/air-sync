import Message from '@/api/models/Message';
import SessionApi from '@/api/SessionApi';
import map from 'lodash/map';
import React from 'react';
import SessionMessage from './SessionMessage';

export interface SessionMessagesProps {
  api: SessionApi;
  sessionId: string;
  messages: Message[];
  timestamp: number;
}

const SessionMessages: React.FC<SessionMessagesProps> = ({
  api,
  sessionId,
  messages,
  timestamp,
}) => {
  return (
    <React.Fragment>
      {map(messages, (message) => (
        <SessionMessage
          key={message.id}
          api={api}
          sessionId={sessionId}
          message={message}
          timestamp={timestamp}
        />
      ))}
    </React.Fragment>
  );
};

export default SessionMessages;
