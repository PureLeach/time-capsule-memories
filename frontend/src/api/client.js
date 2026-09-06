import axios from 'axios';

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_BACKEND_API_URL,
  headers: { 'Content-Type': 'application/json' },
  timeout: 15000,
});

// Collapse axios' error shapes into a plain Error, so callers can render
// `error.message` without knowing which layer failed.
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    const message =
      error.response?.data?.error ||
      error.response?.statusText ||
      error.message ||
      'Request failed';
    const normalized = new Error(message, { cause: error });
    normalized.status = error.response?.status;
    normalized.data = error.response?.data;
    return Promise.reject(normalized);
  }
);

export default apiClient;
