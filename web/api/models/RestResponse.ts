interface RestResponse<T> {
  status: string;
  message?: string;
  data?: T;
  error?: string;
}

export default RestResponse;
