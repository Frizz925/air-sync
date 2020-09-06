import { handleErrorAlert } from '@/utils/Error';
import React from 'react';
import SessionProps from './SessionProps';

export interface CreateSessionProps extends SessionProps {}

const CreateSession: React.FC<CreateSessionProps> = ({ api, connect }) => {
  const createSession = async () => {
    try {
      const { data: sessionId } = await api.createSession();
      connect(sessionId);
    } catch (err) {
      console.error(err);
      handleErrorAlert(err);
    }
  };

  return (
    <div className='text-center'>
      <span
        className='text-blue-500 font-semibold cursor-pointer'
        onClick={createSession}
      >
        Create New Session
      </span>
    </div>
  );
};

export default CreateSession;
