import axios from 'axios';

// 環境変数の末尾に /api がなければ自動で補完する
const getBaseURL = () => {
  const url = process.env.NEXT_PUBLIC_API_URL;
  if (!url) return '/api';
  return url.endsWith('/api') ? url : `${url}/api`;
};

const api = axios.create({
  baseURL: getBaseURL(),
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
