import Axios from 'axios';

const createApiClient = () =>
  Axios.create({
    baseURL: 'http://localhost:8080/api',
  });

export default createApiClient;
