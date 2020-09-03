import Message from './Message';

export interface MessageInserted {
  session_id: string;
  message: Message;
}

export interface MessageDeleted {
  session_id: string;
  message_id: string;
}
