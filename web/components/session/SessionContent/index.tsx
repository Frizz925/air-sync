import Content from '@/api/models/Content';
import Card from '@/components/common/Card';
import React from 'react';

export interface SessionContentProps {
  content: Content;
}

const SessionContent: React.FC<SessionContentProps> = ({ content }) => {
  return (
    <Card className='px-2 py-2 text-sm whitespace-pre-wrap'>
      {content.payload}
    </Card>
  );
};

export default SessionContent;
