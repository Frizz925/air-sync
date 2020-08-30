import Message from '@/api/models/Message';
import Card from '@/components/common/Card';
import IconButton from '@/components/common/IconButton';
import * as Clipboard from '@/utils/Clipboard';
import { faCopy } from '@fortawesome/free-regular-svg-icons';
import moment from 'moment';
import React from 'react';

export interface SessionMessageProps {
  timestamp: number;
  message: Message;
}

const formatTimestamp = (ts: number) => {
  return moment(ts).format('YYYY-MM-DD hh:mm:ss');
};

const formatShortTimestamp = (from: number, ts: number) => {
  return moment(ts).from(from, false);
};

const SessionMessage: React.FC<SessionMessageProps> = ({
  timestamp,
  message,
}) => {
  const handleCopy = () => {
    Clipboard.copy(message.content);
  };

  return (
    <Card className='px-2 py-2 text-sm whitespace-pre-wrap'>
      <div className='flex justify-start items-start'>
        <div>
          <div
            className='text-xs opacity-50 cursor-default'
            title={formatTimestamp(message.created_at)}
          >
            {formatShortTimestamp(timestamp, message.created_at)}
          </div>
          <div>{message.content}</div>
        </div>
        <div className='flex-grow'></div>
        <div>
          <IconButton icon={faCopy} onClick={handleCopy} />
        </div>
      </div>
    </Card>
  );
};

export default SessionMessage;
