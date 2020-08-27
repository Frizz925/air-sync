import { getBaseUrl } from '@/utils/Url';

const createEventSourceClient = (sessionId: string) => {
  const baseUrl = getBaseUrl();
  return new EventSource(`${baseUrl}/sse/sessions/${sessionId}`);
};

export default createEventSourceClient;
