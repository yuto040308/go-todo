// UI モック用のダミーデータ。
// 型は api/openapi.yaml の schemas (Todo) に合わせている。
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
