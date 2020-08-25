import { AxiosInstance } from 'axios';
import Content from './models/Content';
import RestResponse from './models/RestResponse';
import Session from './models/Session';

export default class SessionApi {
  private readonly client: AxiosInstance;

  constructor(client: AxiosInstance) {
    this.client = client;
  }

  public async createSession() {
    const { data } = await this.client.post('/sessions');
    return data as RestResponse<string>;
  }

  public async getSession(id: string) {
    const { data } = await this.client.get(`/sessions/${id}`);
    return data as RestResponse<Session>;
  }

  public async deleteSession(id: string) {
    const { data } = await this.client.delete(`/sessions/${id}`);
    return data as RestResponse<undefined>;
  }

  public async sendMessage(id: string, content: Content) {
    const { data } = await this.client.put(`/sessions/${id}`, content);
    return data as RestResponse<undefined>;
  }
}
