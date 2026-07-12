-- ==========================================
-- RLS 有効化: Data API(PostgREST)経由の情報漏れを防ぐ
-- ==========================================
-- Supabase は public スキーマの全テーブルを Data API で自動公開する。RLS が無効だと
-- publishable key(公開鍵)だけで誰でも全行を読めてしまう(users.password_hash / email 含む)。
-- RLS を有効化し、ポリシーを一切定義しないことで anon(publishable)からの読み取りを 0 件にする。
-- バックエンドは BYPASSRLS を持つ postgres role で接続するため CRUD への影響はなく、
-- keepalive(/rest/v1/todos)も 200 [] を返し続けるので動作は維持される。
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE todos ENABLE ROW LEVEL SECURITY;
ALTER TABLE schema_migrations ENABLE ROW LEVEL SECURITY;