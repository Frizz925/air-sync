import QrImageApi from '@/api/QrImageApi';
import SessionApi from '@/api/SessionApi';
import Dialog from '@/components/common/Dialog';
import { IconProp } from '@fortawesome/fontawesome-svg-core';
import { faTrashAlt } from '@fortawesome/free-regular-svg-icons';
import { faQrcode, faSync } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classNames from 'classnames';
import React, { useState } from 'react';
import styles from './styles.module.css';

export interface SessionActionsProps {
  sessionId: string;
  sessionApi: SessionApi;
  qrImageApi: QrImageApi;
  onReload: () => void;
  onDelete: () => void;
}

const SessionActions: React.FC<SessionActionsProps> = ({
  sessionId,
  sessionApi,
  qrImageApi,
  onReload,
  onDelete,
}) => {
  const [dialogShown, setDialogShown] = useState(false);
  const [qrImageSrc, setQrImageSrc] = useState('');

  const handleQrImage = async () => {
    const href = location.href;
    const queryIdx = href.indexOf('?');
    const link = queryIdx >= 0 ? href.substring(0, queryIdx) : href;
    try {
      const src = await qrImageApi.generate(link);
      setQrImageSrc(src);
      setDialogShown(true);
    } catch (err) {
      console.error(err);
    }
  };

  const handleDelete = async () => {
    try {
      await sessionApi.deleteSession(sessionId);
      onDelete();
    } catch (err) {
      console.error(err);
    }
  };

  const handleClose = () => {
    setDialogShown(false);
  };

  return (
    <React.Fragment>
      <div className='flex flex-row px-2 py-2'>
        <IconButton icon={faQrcode} onClick={handleQrImage} />
        <IconButton icon={faSync} onClick={onReload} />
        <div className='flex-grow'></div>
        <IconButton icon={faTrashAlt} color='red' onClick={handleDelete} />
      </div>
      <Dialog shown={dialogShown} onClose={handleClose}>
        <img src={qrImageSrc} />
      </Dialog>
    </React.Fragment>
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
