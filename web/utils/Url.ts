import { IS_DEV } from './Env';

export const getAttachmentUrl = (id: string) => {
  const baseUrl = getBaseUrl();
  return `${baseUrl}/attachments/${id}`;
};

export const getBaseUrl = (webSocket: boolean = false) => {
  let scheme = getScheme();
  if (webSocket) {
    scheme = scheme === 'https' ? 'wss' : 'ws';
  }
  const host = getHost();
  return `${scheme}://${host}`;
};

export const getScheme = () => {
  if (!process.browser) return 'http';
  const protocol = window.location.protocol;
  return protocol.substring(0, protocol.length - 1);
};

export const getHost = () => {
  if (IS_DEV) {
    return 'localhost:8080';
  } else if (!process.browser) {
    return 'localhost';
  }
  const location = window.location;
  if (location.host) return location.host;
  const hostname = location.hostname;
  const port = location.port;
  return port ? `${hostname}:${port}` : hostname;
};
