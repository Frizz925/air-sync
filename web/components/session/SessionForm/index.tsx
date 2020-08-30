import { DefaultMessage } from '@/api/models/Message';
import SessionApi from '@/api/SessionApi';
import Button from '@/components/common/Button';
import Card from '@/components/common/Card';
import TextBox from '@/components/common/TextBox';
import React, { useState } from 'react';

export interface SessionFormProps {
  api: SessionApi;
  sessionId: string;
}

const SessionForm: React.FC<SessionFormProps> = ({ api, sessionId }) => {
  const [valid, setValid] = useState(false);
  const [processing, setProcessing] = useState(false);
  const [textMessage, setTextMessage] = useState('');

  const resetForm = () => {
    setTextMessage('');
    setValid(false);
  };

  const handleTextChange = (value: string) => {
    const nextValid = !!value;
    if (valid !== nextValid) setValid(nextValid);
    setTextMessage(value);
  };

  const handleSend = async () => {
    if (!valid || processing) return;
    setProcessing(true);
    try {
      const message = DefaultMessage();
      message.content = textMessage;
      await api.sendMessage(sessionId, message);
      resetForm();
    } catch (err) {
      console.error(err);
    } finally {
      setProcessing(false);
    }
  };

  return (
    <Card className='card p-2 space-y-2'>
      <TextBox
        value={textMessage}
        placeholder='Type your message here'
        onChange={handleTextChange}
      />
      <div className='flex flex-row-reverse items-center'>
        <Button
          color='primary'
          className='rounded-full'
          onClick={handleSend}
          disabled={!valid || processing}
        >
          Send
        </Button>
      </div>
    </Card>
  );
};

export default SessionForm;
