import axios from 'axios';
import apiClient from './client';

export async function getPresignedUrl(directory) {
  const { data } = await apiClient.get('/generate-presigned-url', {
    params: { directory },
  });
  return data.presigned_url;
}

export function uploadToS3(presignedUrl, file) {
  return axios.put(presignedUrl, file, {
    headers: { 'Content-Type': file.type },
  });
}

export function createCapsule(payload) {
  return apiClient.post('/capsules', payload);
}
