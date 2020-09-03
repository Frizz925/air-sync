import clsx from 'clsx';
import React from 'react';
import styles from './styles.module.css';

export interface DialogProps {
  shown: boolean;
  onClose: () => void;
}

const Dialog: React.FC<DialogProps> = ({ children, shown, onClose }) => {
  const overlayClasses = clsx(styles.overlay, styles.transition, {
    [styles.shown]: shown,
    [styles.hidden]: !shown,
  });

  const containerClasses = clsx(styles.container, styles.transition, {
    [styles.shown]: shown,
    [styles.hidden]: !shown,
  });

  return (
    <React.Fragment>
      <div className={overlayClasses} onClick={onClose} />
      <div className={containerClasses}>
        <div className='bg-gray-700 rounded-md shadow-lg'>{children}</div>
      </div>
    </React.Fragment>
  );
};

export default Dialog;
