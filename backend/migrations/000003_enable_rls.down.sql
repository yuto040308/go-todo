-- ==========================================
-- RLS 無効化(ロールバック)
-- ==========================================
-- 有効化と逆順で戻す(RLS の ON/OFF は順序非依存だが習慣として)。
ALTER TABLE schema_migrations DISABLE ROW LEVEL SECURITY;
ALTER TABLE todos DISABLE ROW LEVEL SECURITY;
ALTER TABLE users DISABLE ROW LEVEL SECURITY;
