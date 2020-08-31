import Message from './models/Message';
import RestResponse from './models/RestResponse';
import Session from './models/Session';
import RestApi from './RestApi';

export default class SessionApi extends RestApi {
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

  public async sendMessage(id: string, message: Message) {
    const { data } = await this.client.put(`/sessions/${id}`, message);
    return data as RestResponse<undefined>;
  }

  public async deleteMessage(id: string, messageId: string) {
    const { data } = await this.client.delete(`/sessions/${id}/${messageId}`);
    return data as RestResponse<undefined>;
  }
}
