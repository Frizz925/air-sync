interface SessionEvent<T> {
  id: string;
  event: string;
  data: T;
  timestamp: number;
}

export default SessionEvent;
