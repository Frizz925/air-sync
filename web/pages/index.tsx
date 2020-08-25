import SessionApi from '@/api/SessionApi';
import { createApiClient } from '@/clients';
import ConnectSession from '@/components/ConnectSession';
import CreateSession from '@/components/CreateSession';
import { useRouter } from 'next/router';
import React from 'react';

const sessionApi = new SessionApi(createApiClient());

export default function IndexPage() {
  const router = useRouter();

  const connect = (sessionId: string) => {
    router.push(`/s/${sessionId}`);
  };

  return (
    <div className='container container-main'>
      <div className='card max-w-lg mx-auto'>
        <div className='flex flex-col items-center space-y-4 p-4'>
          <h1 className='text-3xl font-semibold'>Air Sync</h1>
          <ConnectSession connect={connect} />
          <CreateSession api={sessionApi} connect={connect} />
        </div>
      </div>
    </div>
  );
}
