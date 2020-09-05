import { IS_BROWSER } from './Env';
import { Observable, ObserverCallback } from './Observable';

const ALERT_OBSERVABLE = 'alertObservable';

export interface AlertMessage {
  type: 'info' | 'error';
  message: string;
  duration?: number;
}

export type ObservableAlert = Observable<AlertMessage>;

export type AlertCallback = ObserverCallback<AlertMessage>;

if (IS_BROWSER && !(ALERT_OBSERVABLE in window)) {
  window[ALERT_OBSERVABLE] = new Observable<AlertMessage>();
}
const observable = IS_BROWSER
  ? (window[ALERT_OBSERVABLE] as ObservableAlert)
  : new Observable<AlertMessage>();

export const subscribe = (callback: AlertCallback) =>
  observable.subscribe(callback);

const alert = (message: AlertMessage) => observable.notify(message);
if (IS_BROWSER) {
  window.alert = (message: any) => alert({ type: 'info', message });
}

export default alert;
