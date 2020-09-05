import alert from './Alert';

export const handleErrorAlert = (err: any) =>
  alert({ type: 'error', message: err });
