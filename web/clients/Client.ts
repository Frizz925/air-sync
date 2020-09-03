import { getBaseUrl } from '@/utils/Url';
import Axios from 'axios';

const createClient = () =>
  Axios.create({
    baseURL: getBaseUrl(),
  });

export default createClient;
