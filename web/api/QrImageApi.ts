import RestApi from './RestApi';

export default class QrImageApi extends RestApi {
  public generate(content: string): Promise<string> {
    return new Promise((resolve, reject) => {
      this.client
        .post('/qr/generate', content, {
          responseType: 'blob',
        })
        .then(({ data }) => {
          const reader = new FileReader();
          reader.onload = () => resolve(reader.result as string);
          reader.onerror = (err) => reject(err);
          reader.readAsDataURL(data);
        });
    });
  }
}
