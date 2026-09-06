import axios from 'axios';
import apiClient from './client';

// The signature covers the key, content type and size limit, so the storage
// server enforces them rather than this client.
export async function getUploadTarget(directory, contentType) {
  const { data } = await apiClient.get('/generate-presigned-url', {
    params: { directory, content_type: contentType },
  });
  return data;
}

export function uploadAttachment({ url, fields }, file) {
  const form = new FormData();
  // Policy fields are signed: send them verbatim, before the file.
  Object.entries(fields).forEach(([key, value]) => form.append(key, value));
  form.append('file', file);

  return axios.post(url, form);
}

export function createCapsule(payload) {
  return apiClient.post('/capsules', payload);
}
