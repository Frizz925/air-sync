import Button from '@/components/common/Button';
import React, { useCallback, useRef } from 'react';
import SessionProps from './SessionProps';

export type ConnectionSessionProps = SessionProps;

const ConnectSession: React.FC<ConnectionSessionProps> = ({ connect }) => {
  const inputRef = useRef<HTMLInputElement>();

  const handleConnect = useCallback(() => {
    const sessionId = inputRef.current.value;
    if (!sessionId) return;
    connect(sessionId);
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
    </div>
  );
};

export default ConnectSession;
