import type { components } from '@/types/api';
import api from '@/lib/api/client';

// 生成型 (types/api.ts) をそのまま使う。API DTO の SoT は api/openapi.yaml。
export type Todo = components['schemas']['Todo'];

// リクエストファイルを用意。GETやDELETEなどリクエストがいらないものは不要
export type CreateTodoRequest = components['schemas']['CreateTodoRequest'];
export type UpdateTodoRequest = components['schemas']['UpdateTodoRequest'];

// 1件のTODOを取得
export const getTodo = async (id: string): Promise<Todo> => {
  const res = await api.get<Todo>(`/todos/${id}`);
  return res.data;
};

// 複数件のTODOを取得
export const listTodos = async (): Promise<Todo[]> => {
  const res = await api.get<Todo[]>('/todos');
  return res.data;
};

// TODOを作成
export const createTodo = async (body: CreateTodoRequest): Promise<Todo> => {
  const res = await api.post<Todo>('/todos', body);
  return res.data;
};

// TODOを更新
export const updateTodo = async (id: string, body: UpdateTodoRequest): Promise<Todo> => {
  const res = await api.put<Todo>(`/todos/${id}`, body);
  return res.data;
};

// TODOを削除
export const deleteTodo = async (id: string): Promise<void> => {
  await api.delete(`/todos/${id}`);
};
