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
    <div>
      <div className={overlayClasses} onClick={onClose}></div>
      <div className={containerClasses}>
        <div className={styles.content}>{children}</div>
      </div>
    </div>
  );
};

export default Dialog;
