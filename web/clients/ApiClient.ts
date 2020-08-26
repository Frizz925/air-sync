import { getBaseUrl } from '@/utils/Url';
import Axios from 'axios';

const createApiClient = () => {
  const baseUrl = getBaseUrl();
  return Axios.create({ baseURL: `${baseUrl}/api` });
};

export default createApiClient;
