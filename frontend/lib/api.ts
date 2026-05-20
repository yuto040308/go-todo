import axios from 'axios';

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || '/api', // 環境変数から動的に読み込み（未設定時はローカルのプロキシ /api）
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

