import Content from '@/api/models/Content';
import SessionApi from '@/api/SessionApi';
import React, { useRef, useState } from 'react';
import styles from './styles.module.css';

export interface SessionFormProps {
  api: SessionApi;
  sessionId: string;
}

const SessionForm: React.FC<SessionFormProps> = ({ api, sessionId }) => {
  const [valid, setValid] = useState(false);
  const [processing, setProcessing] = useState(false);
  const textRef = useRef<HTMLTextAreaElement>();

  const resetForm = () => {
    textRef.current.value = '';
    setValid(false);
  };

  const handleTextArea = (evt: React.ChangeEvent<HTMLTextAreaElement>) => {
    const value = evt.target.value;
    setValid(!!value);
  };

  const handleSend = async () => {
    setProcessing(true);
    try {
      const content: Content = {
        type: 'text',
        mime: 'text/plain',
        payload: textRef.current.value,
      };
      await api.sendMessage(sessionId, content);
      resetForm();
    } catch (err) {
      console.error(err);
    } finally {
      setProcessing(false);
    }
  };

  return (
    <div className='card p-2'>
      <div>
        <textarea
          ref={textRef}
          className={styles.textarea}
          rows={1}
          placeholder='Type your message here'
          disabled={processing}
          onChange={handleTextArea}
        />
      </div>
      <div className='flex flex-row-reverse items-center'>
        <button
          className='btn btn-primary rounded-full'
          onClick={handleSend}
          disabled={!valid}
        >
          Send
        </button>
      </div>
    </div>
  );
};

export default SessionForm;
