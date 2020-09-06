import Message from '@/api/models/Message';
import map from 'lodash/map';
import React from 'react';
import SessionMessage from './SessionMessage';

export interface SessionMessagesProps {
  messages: Message[];
  timestamp: number;
  onDelete: (message: Message) => void;
}

const SessionMessages: React.FC<SessionMessagesProps> = ({
  messages,
  timestamp,
  onDelete,
}) => {
  return (
    <React.Fragment>
      {map(messages, (message) => (
        <SessionMessage
          key={message.id}
          message={message}
          timestamp={timestamp}
          onDelete={() => onDelete(message)}
        />
      ))}
    </React.Fragment>
  );
};

export default SessionMessages;
