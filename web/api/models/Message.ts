interface Message {
  id: string;
  sensitive: boolean;
  body?: string;
  attachment_id?: string;
  attachment_type?: string;
  attachment_name?: string;
  created_at: number;
}

export const DefaultMessage = (): Partial<Message> => ({
  sensitive: false,
  created_at: new Date().getTime() / 1000,
});

export default Message;
