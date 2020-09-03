import Message from '@/api/models/Message';
import { getAttachmentUrl } from '@/utils/Url';
import clsx from 'clsx';
import React, { useCallback, useMemo, useState } from 'react';
import styles from './styles.module.css';

export interface MessageContentProps {
  message: Message;
}

const MessageContent: React.FC<MessageContentProps> = ({ message }) => {
  const { body, attachment_id, sensitive } = message;
  const [revealed, setRevealed] = useState(false);
  const toggleReveal = useCallback(() => setRevealed(!revealed), [revealed]);
  const textCls = useMemo(
    () =>
      clsx(styles.text, {
        [styles.sensitive]: sensitive,
        [styles.revealed]: revealed,
      }),
    [sensitive, revealed]
  );
  return (
    <div className='space-y-2'>
      {body && (
        <div className='px-2'>
          <span className={textCls} onClick={toggleReveal}>
            {body}
          </span>
        </div>
      )}
      {attachment_id && <MessageAttachment message={message} />}
    </div>
  );
};

const MessageAttachment: React.FC<{ message: Message }> = ({
  message: { attachment_id: id, attachment_type: type, sensitive },
}) => {
  const [revealed, setRevealed] = useState(false);
  const url = useMemo(() => getAttachmentUrl(id), [id]);
  const toggleReveal = useCallback(() => setRevealed(!revealed), [revealed]);
  const imgCls = useMemo(() => {
    const cond = {
      [styles.sensitive]: sensitive,
      [styles.revealed]: revealed,
      [styles.unrevealed]: !revealed,
    };
    return {
      container: clsx(styles.imageContainer, cond),
      image: clsx(styles.image, cond),
    };
  }, [sensitive, revealed]);
  return type === 'image' ? (
    <div className='flex justify-center'>
      <div className={imgCls.container} onClick={toggleReveal}>
        <img className={imgCls.image} src={url} />
      </div>
    </div>
  ) : null;
};

export default MessageContent;
