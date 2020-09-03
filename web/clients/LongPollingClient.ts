import { getBaseUrl } from '@/utils/Url';
import Axios from 'axios';

const createLongPollingClient = () => {
  const baseUrl = getBaseUrl();
  return Axios.create({
    baseURL: `${baseUrl}/lp`,
  });
};

export default createLongPollingClient;
