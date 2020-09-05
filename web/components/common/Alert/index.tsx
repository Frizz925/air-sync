import { AlertMessage } from '@/utils/Alert';
import clsx from 'clsx';
import React from 'react';
import styles from './styles.module.css';

export interface AlertProps {
  shown: boolean;
  alert: AlertMessage;
}

const Alert: React.FC<AlertProps> = ({ shown, alert }) => {
  const outerCls = clsx(styles.outer, shown && styles.shown);
  const containerCls = clsx(styles.container);
  const cardCls = clsx(styles.card, styles[alert.type]);
  return (
    <div className={outerCls}>
      <div className={containerCls}>
        <div className={cardCls}>{alert.message}</div>
      </div>
    </div>
  );
};

export default Alert;
