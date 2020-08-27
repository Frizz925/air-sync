import { getEnvBool } from './Env';

export class NotificationHelper {
  private notificationAllowed = false;

  public initialize() {
    if (!process.browser) return;
    if (!getEnvBool('NOTIFICATION_ENABLED')) return;
    Notification.requestPermission().then(
      (result) => {
        this.notificationAllowed = result === 'granted';
      },
      (err) => console.error(err)
    );
  }

  public notify(title: string, body: string) {
    if (!this.notificationAllowed) return;
    new Notification(title, { body });
  }
}
