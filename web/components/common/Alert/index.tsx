import { AlertMessage } from '@/utils/Alert';
import {
  faExclamationTriangle,
  faInfoCircle,
} from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import clsx from 'clsx';
import React from 'react';
import styles from './styles.module.css';

export interface AlertProps {
  shown: boolean;
  alert: AlertMessage;
}

const Alert: React.FC<AlertProps> = ({ shown, alert }) => {
  let icon = faInfoCircle;
  switch (alert.type) {
    case 'error':
      icon = faExclamationTriangle;
      break;
  }

  const outerCls = clsx(styles.outer, shown && styles.shown);
  const containerCls = clsx(styles.container);
  const cardCls = clsx(styles.card, styles[alert.type]);
  return (
    <div className={outerCls}>
      <div className={containerCls}>
        <div className={cardCls}>
          <div>
            <FontAwesomeIcon icon={icon} />
          </div>
          <div className={styles.message}>{alert.message}</div>
        </div>
      </div>
    </div>
  );
};

export default Alert;
