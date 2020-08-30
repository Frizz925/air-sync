import Message from '@/api/models/Message';
import Card from '@/components/common/Card';
import moment from 'moment';
import React from 'react';

export interface SessionMessageProps {
  timestamp: number;
  message: Message;
}

const formatTimestamp = (from: number, ts: number) => {
  return moment(ts).from(from, false);
};

const SessionMessage: React.FC<SessionMessageProps> = ({
  timestamp,
  message,
}) => {
  return (
    <Card className='px-2 py-2 text-sm whitespace-pre-wrap'>
      <div className='text-xs opacity-50'>
        {formatTimestamp(timestamp, message.created_at)}
      </div>
      <div>{message.content}</div>
    </Card>
  );
};

export default SessionMessage;
