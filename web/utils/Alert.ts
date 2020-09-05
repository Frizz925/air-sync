import { Observable, ObserverCallback } from './Observable';

const ALERT_OBSERVABLE = 'alertObservable';

export interface AlertMessage {
  type: 'info' | 'error';
  message: string;
}

if (!(ALERT_OBSERVABLE in window)) {
  window[ALERT_OBSERVABLE] = new Observable<AlertMessage>();
}
const observable = window[ALERT_OBSERVABLE] as Observable<AlertMessage>;

export const subscribe = (callback: ObserverCallback<AlertMessage>) =>
  observable.subscribe(callback);

const alert = (message: AlertMessage) => observable.notify(message);
export default alert;
