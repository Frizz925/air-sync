import Message from '@/api/models/Message';
import Card from '@/components/common/Card';
import IconButton from '@/components/common/IconButton';
import alert from '@/utils/Alert';
import * as Clipboard from '@/utils/Clipboard';
import { formatShortTimestamp, formatTimestamp } from '@/utils/Time';
import { getAttachmentUrl } from '@/utils/Url';
import {
  faCloudDownloadAlt as faDownload,
  faCopy,
  faTrashAlt,
} from '@fortawesome/free-solid-svg-icons';
import React from 'react';
import MessageContent from '../MessageContent';
import styles from './styles.module.css';

export interface SessionMessageProps {
  message: Message;
  timestamp: number;
  onDelete: () => void;
}

const SessionMessage: React.FC<SessionMessageProps> = ({
  message,
  timestamp,
  onDelete,
}) => {
  const {
    attachment_id: attachmentId,
    attachment_name: attachmentName,
    created_at: createdAt,
  } = message;
  const handleCopy = () => {
    Clipboard.copy(message.body);
    alert('Message copied to clipboard');
  };

  const attachmentUrl = attachmentId
    ? getAttachmentUrl(attachmentId)
    : undefined;

  return (
    <Card className='text-sm whitespace-pre-wrap'>
      <div className='pt-2 space-y-2'>
        <div
          className='px-2 text-xs opacity-50 cursor-default'
          title={formatTimestamp(createdAt)}
        >
          {formatShortTimestamp(createdAt, timestamp)}
        </div>
        <MessageContent message={message} />
      </div>
      <div className='flex items-center px-1 py-1'>
        {message.body && (
          <div>
            <IconButton icon={faCopy} onClick={handleCopy} />
          </div>
        )}
        {attachmentUrl && (
          <div className='flex items-center overflow-hidden'>
            <a
              href={attachmentUrl}
              title={attachmentName}
              download={attachmentName}
              target='_blank'
              rel='noreferrer'
            >
              <IconButton icon={faDownload} />
            </a>
            <div className={styles.attachmentName}>
              <a
                href={attachmentUrl}
                title={attachmentName}
                target='_blank'
                rel='noreferrer'
              >
                {attachmentName}
              </a>
            </div>
          </div>
        )}
        <div className='flex-grow'></div>
        <div>
          <IconButton icon={faTrashAlt} color='red' onClick={onDelete} />
        </div>
      </div>
    </Card>
  );
};

export default SessionMessage;
