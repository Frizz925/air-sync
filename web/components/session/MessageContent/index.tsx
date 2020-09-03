import Message from '@/api/models/Message';
import { getAttachmentUrl } from '@/utils/Url';
import React from 'react';

export interface MessageContentProps {
  message: Message;
}

const MessageContent: React.FC<MessageContentProps> = ({ message }) => {
  return (
    <div className='space-y-2'>
      {message.body && <div className='px-2'>{message.body}</div>}
      {message.attachment_id && <MessageAttachment message={message} />}
    </div>
  );
};

const MessageAttachment: React.FC<{ message: Message }> = ({
  message: { attachment_id: id, attachment_type: type },
}) => {
  const url = getAttachmentUrl(id);
  return type === 'image' ? (
    <div className='flex justify-center'>
      <img src={url} />
    </div>
  ) : null;
};

export default MessageContent;
