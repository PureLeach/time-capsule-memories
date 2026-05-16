import apiClient from './client';

export function submitFeedback(payload) {
  return apiClient.post('/feedback', payload);
}
