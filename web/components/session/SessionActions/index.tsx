import QrImageApi from '@/api/QrImageApi';
import SessionApi from '@/api/SessionApi';
import Dialog from '@/components/common/Dialog';
import IconButton from '@/components/common/IconButton';
import * as Clipboard from '@/utils/Clipboard';
import { handleErrorAlert } from '@/utils/Error';
import {
  faCopy,
  faQrcode,
  faSync,
  faTrashAlt,
} from '@fortawesome/free-solid-svg-icons';
import React, { useState } from 'react';

const getCurrentUrl = () => {
  const href = location.href;
  const queryIdx = href.indexOf('?');
  return queryIdx >= 0 ? href.substring(0, queryIdx) : href;
};

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
    try {
      const src = await qrImageApi.generate(getCurrentUrl());
      setQrImageSrc(src);
      setDialogShown(true);
    } catch (err) {
      console.error(err);
      handleErrorAlert(err);
    }
  };

  const handleCopy = () => Clipboard.copy(getCurrentUrl());

  const handleDelete = async () => {
    try {
      await sessionApi.deleteSession(sessionId);
      onDelete();
    } catch (err) {
      console.error(err);
      handleErrorAlert(err);
    }
  };

  const handleClose = () => {
    setDialogShown(false);
  };

  return (
    <React.Fragment>
      <div className='flex flex-row px-1 py-2'>
        <IconButton
          icon={faQrcode}
          color='blue'
          title='Generate QR code'
          onClick={handleQrImage}
        />
        <IconButton
          icon={faCopy}
          color='blue'
          title='Copy session URL'
          onClick={handleCopy}
        />
        <IconButton
          icon={faSync}
          color='blue'
          title='Reload'
          onClick={onReload}
        />
        <div className='flex-grow'></div>
        <IconButton
          icon={faTrashAlt}
          color='red'
          title='Delete session'
          onClick={handleDelete}
        />
      </div>
      <Dialog shown={dialogShown} onClose={handleClose}>
        <img src={qrImageSrc} />
      </Dialog>
    </React.Fragment>
  );
};

export default SessionActions;
