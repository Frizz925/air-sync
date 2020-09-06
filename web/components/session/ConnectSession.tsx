import RestResponse from '@/api/models/RestResponse';
import Button from '@/components/common/Button';
import { handleErrorAlert } from '@/utils/Error';
import { AxiosResponse } from 'axios';
import React, { useCallback, useState } from 'react';
import SessionProps from './SessionProps';

export interface ConnectionSessionProps extends SessionProps {}

const ConnectSession: React.FC<ConnectionSessionProps> = ({ api, connect }) => {
  const [value, setValue] = useState('');

  const handleConnect = useCallback(async () => {
    const sessionId = value;
    if (!sessionId) return;
    try {
      await api.getSession(sessionId);
      connect(sessionId);
    } catch (err) {
      if (err.response) {
        const res = (err.response as AxiosResponse).data as RestResponse<{}>;
        const message = res.error || res.message || "Can't connect to session";
        handleErrorAlert(message);
      } else {
        console.error(err);
        handleErrorAlert(err);
      }
    }
  }, [value, connect]);

  return (
    <div className='flex flex-col items-center w-full space-y-2'>
      <input
        type='text'
        className='bg-gray-700 w-full px-4 py-2 rounded-full outline-none text-center'
        placeholder='Enter the session ID'
        value={value}
        onChange={(e) => setValue(e.target.value)}
      />
      <Button
        color='primary'
        className='rounded-full'
        disabled={!value}
        onClick={handleConnect}
      >
        Connect
      </Button>
    </div>
  );
};

export default ConnectSession;
