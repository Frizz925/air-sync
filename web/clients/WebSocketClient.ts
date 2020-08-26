import { getBaseUrl } from '@/utils/Url';

const createClient = (sessionId: string) => {
  const baseUrl = getBaseUrl(true);
  return new WebSocket(`${baseUrl}/ws/sessions/${sessionId}`);
};

export default createClient;
