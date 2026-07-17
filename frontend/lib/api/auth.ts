import api from '@/lib/api/client';
import type { components } from '@/types/api';

// 生成型 (types/api.ts) をそのまま使う。API DTO の SoT は api/openapi.yaml。
export type User = components['schemas']['User'];
export type SignupRequest = components['schemas']['SignupRequest'];
export type LoginRequest = components['schemas']['LoginRequest'];
export type LoginResponse = components['schemas']['LoginResponse'];

// 新規登録。仕様上 201 で User のみ返る (token は返らない) ため、
// 呼び出し側で続けて login する 2 段構成になる。
export const signup = async (body: SignupRequest): Promise<User> => {
  const res = await api.post<User>('/auth/signup', body);
  return res.data;
};

// ログイン。JWT (token) + user を返す。token の保管は AuthContext 側の責務。
export const login = async (body: LoginRequest): Promise<LoginResponse> => {
  const res = await api.post<LoginResponse>('/auth/login', body);
  return res.data;
};

// ログアウト。現状サーバ側は 204 を返すのみ (トークン破棄はクライアント責務)。
// interceptor が Bearer を載せるため auth 必須でも通る。
export const logout = async (): Promise<void> => {
  await api.post('/auth/logout');
};

// 現在ログイン中のユーザー情報。認証状態の判定に使う。
export const getMe = async (): Promise<User> => {
  const res = await api.get<User>('/auth/me');
  return res.data;
};
