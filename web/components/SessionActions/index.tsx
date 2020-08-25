import SessionApi from '@/api/SessionApi';
import { IconProp } from '@fortawesome/fontawesome-svg-core';
import { faTrashAlt } from '@fortawesome/free-regular-svg-icons';
import { faQrcode } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classNames from 'classnames';
import React from 'react';
import styles from './styles.module.css';

export interface SessionActionsProps {
  api: SessionApi;
  sessionId: string;
}

const SessionActions: React.FC<SessionActionsProps> = ({ api, sessionId }) => {
  const handleQrCode = () => {};

  const handleDelete = async () => {
    try {
      await api.deleteSession(sessionId);
    } catch (err) {
      console.error(err);
    }
  };
  return (
    <div className='flex flex-row'>
      <IconButton icon={faQrcode} onClick={handleQrCode} />
      <IconButton icon={faTrashAlt} color='red' onClick={handleDelete} />
    </div>
  );
};

interface IconButtonProps {
  color?: string;
  icon: IconProp;
  onClick: () => void;
}

const IconButton: React.FC<IconButtonProps> = ({ color, icon, onClick }) => {
  const classes = color
    ? classNames(styles.action, styles[color])
    : styles.action;
  return (
    <button className={classes} onClick={onClick}>
      <FontAwesomeIcon icon={icon} />
    </button>
  );
};

export default SessionActions;
