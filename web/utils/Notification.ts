const notificationEnabled =
  process.env.NEXT_PUBLIC_NOTIFICATION_ENABLED === 'true';

export class NotificationHelper {
  private notificationAllowed = false;

  public initialize() {
    if (!process.browser) return;
    if (!notificationEnabled) return;
    if (!('Notification' in window)) return;
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
