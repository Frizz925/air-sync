import classNames from 'classnames';
import React from 'react';

const Card: React.FC<{ className?: string }> = ({ children, className }) => {
  const classes = classNames(
    'bg-gray-800 text-gray-300 rounded-lg shadow-xl mx-auto overflow-hidden',
    className
  );
  return <div className={classes}>{children}</div>;
};

export default Card;
