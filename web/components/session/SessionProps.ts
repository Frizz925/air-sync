import SessionApi from '@/api/SessionApi';

export interface SessionProps {
  api: SessionApi;
  connect: (sessionId: string) => void;
}

export default SessionProps;
