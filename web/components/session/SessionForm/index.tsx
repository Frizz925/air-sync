import Content from '@/api/models/Content';
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
  const [textContent, setTextContent] = useState('');

  const resetForm = () => {
    setTextContent('');
    setValid(false);
  };

  const handleTextChange = (value: string) => {
    const nextValid = !!value;
    if (valid !== nextValid) setValid(nextValid);
    setTextContent(value);
  };

  const handleSend = async () => {
    if (!valid || processing) return;
    setProcessing(true);
    try {
      const content: Content = {
        type: 'text',
        mime: 'text/plain',
        payload: textContent,
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
    <Card className='card p-2 space-y-2'>
      <TextBox
        value={textContent}
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
