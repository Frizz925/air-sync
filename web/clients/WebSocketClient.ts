const baseUrl = 'ws://localhost:8080/ws';

const createClient = (sessionId: string) =>
  new WebSocket(`${baseUrl}/sessions/${sessionId}`);

export default createClient;
