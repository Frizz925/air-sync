interface Attachment {
  id: string;
  type: string;
  mime: string;
  name: string;
  created_at: number;
}

export interface CreateAttachment {
  type: string;
  mime: string;
  name: string;
}

export default Attachment;
