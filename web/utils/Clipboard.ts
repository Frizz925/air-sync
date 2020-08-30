export const copy = (content: string) => {
  const input = document.createElement('input');
  document.body.appendChild(input);

  input.value = content;
  input.select();
  input.setSelectionRange(0, 99999);
  document.execCommand('copy');

  document.body.removeChild(input);
};
