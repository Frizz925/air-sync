interface SessionEvent<T> {
  event: string;
  data: T;
  timestamp: number;
}

export default SessionEvent;
