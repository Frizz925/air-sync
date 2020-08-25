interface Content {
  type: string;
  mime: string;
  payload: string;
}

export const DefaultContent = (): Content => ({
  type: 'text',
  mime: 'text/plain',
  payload: '',
});

export default Content;
