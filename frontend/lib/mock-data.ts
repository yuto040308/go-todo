// UI モック用のダミーデータ。
// 型は api/openapi.yaml の schemas (Todo / User) に合わせている。
// チケット7 の実装で実 API + 生成型 (types/api.ts) に差し替える想定。

export type Todo = {
  id: string;
  user_id: string;
  title: string;
  description: string | null;
  is_completed: boolean;
  created_at: string;
  updated_at: string;
};

export type User = {
  id: string;
  user_name: string;
  email: string;
};

export const mockUser: User = {
  id: '11111111-1111-1111-1111-111111111111',
  user_name: 'yuto',
  email: 'yuto@example.com',
};

export const mockTodos: Todo[] = [
  {
    id: 'a0000000-0000-0000-0000-000000000001',
    user_id: mockUser.id,
    title: 'OpenAPI のスキーマを確認する',
    description: 'Todo / User の型をモックに反映する',
    is_completed: true,
    created_at: '2026-06-25T09:00:00Z',
    updated_at: '2026-06-25T10:30:00Z',
  },
  {
    id: 'a0000000-0000-0000-0000-000000000002',
    user_id: mockUser.id,
    title: 'ログイン画面のモックを作る',
    description: 'email + password、サインアップ導線つき',
    is_completed: false,
    created_at: '2026-06-26T11:00:00Z',
    updated_at: '2026-06-26T11:00:00Z',
  },
  {
    id: 'a0000000-0000-0000-0000-000000000003',
    user_id: mockUser.id,
    title: '一覧画面のレイアウトを詰める',
    description: null,
    is_completed: false,
    created_at: '2026-06-27T14:20:00Z',
    updated_at: '2026-06-27T14:20:00Z',
  },
  {
    id: 'a0000000-0000-0000-0000-000000000004',
    user_id: mockUser.id,
    title: 'スマホ表示で崩れないか確認する',
    description: 'Tailwind のブレークポイントで最低限の対応',
    is_completed: false,
    created_at: '2026-06-28T08:05:00Z',
    updated_at: '2026-06-28T08:05:00Z',
  },
];
