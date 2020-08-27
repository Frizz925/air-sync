import React from 'react';
import SessionProps from './SessionProps';

export interface CreateSessionProps extends SessionProps {}

const CreateSession: React.FC<CreateSessionProps> = ({ api, connect }) => {
  const createSession = async () => {
    const { data: sessionId } = await api.createSession();
    connect(sessionId);
  };

  return (
    <div className='text-center'>
      <a
        href='#'
        className='text-blue-500 font-semibold'
        onClick={createSession}
      >
        Create New Session
      </a>
    </div>
  );
};

export default CreateSession;
