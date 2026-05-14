import axios from 'axios';

const api = axios.create({
  baseURL: '/api', // Go バックエンドの URL
  headers: {
    'Content-Type': 'application/json',
  },
});

// 全てのリクエストにトークンを乗せる
api.interceptors.request.use((config) => {
  // localStorageはブラウザ専用のため、SSR環境(Server Components)で参照しないようにガードする
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
  if (token) {
    // ヘッダーに "Bearer <トークン>" をセット
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default api;
