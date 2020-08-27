import RestResponse from '@/api/models/RestResponse';
import Button from '@/components/common/Button';
import { AxiosResponse } from 'axios';
import React, { useCallback, useRef, useState } from 'react';
import SessionProps from './SessionProps';

export interface ConnectionSessionProps extends SessionProps {}

const ConnectSession: React.FC<ConnectionSessionProps> = ({ api, connect }) => {
  const [error, setError] = useState('');
  const inputRef = useRef<HTMLInputElement>();

  const handleConnect = useCallback(async () => {
    const sessionId = inputRef.current.value;
    if (!sessionId) return;
    try {
      await api.getSession(sessionId);
      connect(sessionId);
    } catch (err) {
      let message = "Can't connect to session";
      if (err.response) {
        const res = (err.response as AxiosResponse).data as RestResponse<{}>;
        message = res.error || res.message || message;
      } else {
        console.error(err);
      }
      setError(message);
    }
  }, [inputRef, connect]);

  return (
    <div className='flex flex-col items-center w-full space-y-2'>
      <input
        className='bg-gray-700 w-full px-4 py-2 rounded-full outline-none text-center'
        placeholder='Enter the session ID'
        type='text'
        ref={inputRef}
      />
      <Button color='primary' className='rounded-full' onClick={handleConnect}>
        Connect
      </Button>
      <div className='text-red-700 text-sm font-semibold'>{error}</div>
    </div>
  );
};

export default ConnectSession;
