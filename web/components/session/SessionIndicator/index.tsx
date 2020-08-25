import ConnectionState from '@/components/models/ConnectionState';
import classNames from 'classnames';
import React from 'react';

export interface SessionIndicatorProps {
  sessionId: string;
  connectionState: ConnectionState;
}

const SessionIndicator: React.FC<SessionIndicatorProps> = ({
  sessionId,
  connectionState,
}) => {
  const indicatorClasses = classNames('p-1 rounded-full', {
    'bg-red-700': connectionState === ConnectionState.DISCONNECTED,
    'bg-yellow-300': connectionState === ConnectionState.CONNECTING,
    'bg-green-600': connectionState === ConnectionState.CONNECTED,
  });

  return (
    <div className='flex flex-row items-center'>
      <div className='mr-4'>
        <div className={indicatorClasses}></div>
      </div>
      <div className='flex flex-col flex-grow'>
        <span className='font-semibold'>Session ID</span>
        <span className='text-sm'>{sessionId}</span>
      </div>
    </div>
  );
};

export default SessionIndicator;
