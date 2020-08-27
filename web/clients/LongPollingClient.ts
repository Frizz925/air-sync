import { getBaseUrl } from '@/utils/Url';
import Axios from 'axios';

const createLongPollingClient = () => {
  const baseUrl = getBaseUrl();
  return Axios.create({
    baseURL: `${baseUrl}/lp`,
    timeout: 30000,
  });
};

export default createLongPollingClient;
