import Message from '@/api/models/Message';
import SessionApi from '@/api/SessionApi';
import Card from '@/components/common/Card';
import IconButton from '@/components/common/IconButton';
import * as Clipboard from '@/utils/Clipboard';
import { faCopy, faTrashAlt } from '@fortawesome/free-regular-svg-icons';
import moment from 'moment';
import React from 'react';

export interface SessionMessageProps {
  api: SessionApi;
  sessionId: string;
  message: Message;
  timestamp: number;
}

const formatTimestamp = (ts: number) => {
  return moment(ts).format('YYYY-MM-DD hh:mm:ss');
};

const formatShortTimestamp = (from: number, ts: number) => {
  return moment(ts).from(from, false);
};

const SessionMessage: React.FC<SessionMessageProps> = ({
  api,
  sessionId,
  message,
  timestamp,
}) => {
  const handleCopy = () => Clipboard.copy(message.content);

  const handleDelete = async () => {
    try {
      await api.deleteMessage(sessionId, message.id);
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <Card className='text-sm whitespace-pre-wrap'>
      <div className='px-2 pt-2'>
        <div
          className='text-xs opacity-50 cursor-default'
          title={formatTimestamp(message.created_at)}
        >
          {formatShortTimestamp(timestamp, message.created_at)}
        </div>
        <div>{message.content}</div>
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
