import each from 'lodash/each';

export type ObserverMap<T> = {
  [key: number]: Observer<T>;
};

export type ObserverCallback<T> = (value: T, error?: any) => void;

export class Observable<T> {
  private readonly observers: ObserverMap<T> = {};

  private nextId = 1;

  public subscribe(callback: ObserverCallback<T>): Observer<T> {
    const id = this.nextId++;
    const observer = new Observer(id, this, callback);
    this.observers[id] = observer;
    return observer;
  }

  public notify(value: T, error?: any) {
    each(this.observers, (observer) => observer.notify(value, error));
  }

  public unsubscribe(id: number) {
    delete this.observers[id];
  }
}

export class Observer<T> {
  private readonly id: number;
  private readonly observable: Observable<T>;
  private readonly callback: ObserverCallback<T>;

  constructor(
    id: number,
    observable: Observable<T>,
    callback: ObserverCallback<T>
  ) {
    this.id = id;
    this.observable = observable;
    this.callback = callback;
  }

  public notify(value: T, error?: any) {
    this.callback(value, error);
  }

  public unsubscribe() {
    this.observable.unsubscribe(this.id);
  }
}

const create = <T>() => new Observable<T>();
export default create;
