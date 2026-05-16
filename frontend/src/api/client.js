import axios from 'axios';

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_BACKEND_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    const message =
      error.response?.data?.message ||
      error.response?.statusText ||
      error.message ||
      'Request failed';
    const normalized = new Error(message);
    normalized.status = error.response?.status;
    normalized.data = error.response?.data;
    normalized.cause = error;
    return Promise.reject(normalized);
  }
);

export default apiClient;
