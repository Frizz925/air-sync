import Content from '@/api/models/Content';
import React from 'react';

export interface SessionContentProps {
  content: Content;
}

const SessionContent: React.FC<SessionContentProps> = ({ content }) => {
  const visible = !!content.payload;
  return <div className='card px-4 py-2 text-sm'>{content.payload}</div>;
};

export default SessionContent;
