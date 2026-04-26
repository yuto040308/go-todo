import axios from 'axios';

const api = axios.create({
  baseURL: '/api', // Go バックエンドの URL
  headers: {
    'Content-Type': 'application/json',
  },
});

// 全てのリクエストにトークンを乗せる
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    // ヘッダーに "Bearer <トークン>" をセット
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default api;