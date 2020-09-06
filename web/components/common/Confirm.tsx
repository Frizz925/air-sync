import React from 'react';
import Button, { ButtonColor } from './Button';
import Card from './Card';
import Dialog from './Dialog';

export interface ConfirmProps {
  shown: boolean;
  message: string;
  confirmLabel?: string;
  cancelLabel?: string;
  confirmColor?: ButtonColor;
  cancelColor?: ButtonColor;
  onConfirm: () => void;
  onCancel?: () => void;
  onClose: () => void;
}

const Confirm: React.FC<ConfirmProps> = ({
  shown,
  message,
  confirmLabel,
  cancelLabel,
  confirmColor,
  cancelColor,
  onConfirm,
  onCancel,
  onClose,
}) => {
  return (
    <Dialog shown={shown} onClose={onClose}>
      <Card className='px-2 py-4 w-full max-w-md'>
        <div className='text-center'>
          <div className='mb-4'>{message}</div>
          <div className='space-x-2'>
            <Button rounded color={cancelColor} onClick={onCancel || onClose}>
              {cancelLabel || 'Cancel'}
            </Button>
            <Button rounded color={confirmColor} onClick={onConfirm}>
              {confirmLabel || 'OK'}
            </Button>
          </div>
        </div>
      </Card>
    </Dialog>
  );
};

export default Confirm;
