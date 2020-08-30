import Message from './Message';

interface Session {
  id: string;
  messages: Message[];
  created_at: number;
}

export default Session;
