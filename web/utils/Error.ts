import { AxiosError } from 'axios';
import alert from './Alert';

export const handleErrorAlert = (err: any) => {
  let message = 'Unexpected error occured';
  if (typeof err === 'string') {
    message = err;
  } else if (typeof err === 'object') {
    if (err.isAxiosError) {
      const axErr = err as AxiosError;
      const res = axErr.response;
      if (res) {
        const data = res.data;
        if (typeof data === 'object') {
          if (data.error) message = data.error;
          else if (data.message) message = data.message;
          else message = axErr.message;
        } else if (typeof data === 'string') {
          message = data;
        } else {
          message = axErr.message;
        }
      } else {
        message = axErr.message;
      }
    } else if (err instanceof Error) {
      message = err.message;
    }
  }
  alert({ type: 'error', message, duration: 5000 });
};
