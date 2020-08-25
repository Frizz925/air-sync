import { AxiosInstance } from 'axios';

export default abstract class RestApi {
  protected readonly client: AxiosInstance;

  constructor(client: AxiosInstance) {
    this.client = client;
  }
}
