import Message from '@/api/models/Message';
import SessionApi from '@/api/SessionApi';
import Card from '@/components/common/Card';
import IconButton from '@/components/common/IconButton';
import * as Clipboard from '@/utils/Clipboard';
import { formatShortTimestamp, formatTimestamp } from '@/utils/Time';
import { faCopy, faTrashAlt } from '@fortawesome/free-regular-svg-icons';
import React from 'react';
import MessageContent from '../MessageContent';

export interface SessionMessageProps {
  api: SessionApi;
  sessionId: string;
  message: Message;
  timestamp: number;
}

const SessionMessage: React.FC<SessionMessageProps> = ({
  api,
  sessionId,
  message,
  timestamp,
}) => {
  const handleCopy = () => Clipboard.copy(message.body);

  const handleDelete = async () => {
    try {
      await api.deleteMessage(sessionId, message.id);
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <Card className='text-sm whitespace-pre-wrap'>
      <div className='pt-2 space-y-2'>
        <div
          className='px-2 text-xs opacity-50 cursor-default'
          title={formatTimestamp(message.created_at)}
        >
          {formatShortTimestamp(message.created_at, timestamp)}
        </div>
        <MessageContent message={message} />
      </div>
      <div className='flex justify-start items-stretch px-1 py-1'>
        <IconButton icon={faCopy} onClick={handleCopy} />
        <div className='flex-grow'></div>
        <IconButton icon={faTrashAlt} color='red' onClick={handleDelete} />
      </div>
    </Card>
  );
};

export default SessionMessage;
