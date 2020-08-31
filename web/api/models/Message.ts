interface Message {
  id: string;
  type: string;
  mime: string;
  content: string;
  created_at: number;
}

export const DefaultMessage = (): Partial<Message> => ({
  type: 'text',
  mime: 'text/plain',
  content: '',
  created_at: new Date().getTime(),
});

export default Message;
