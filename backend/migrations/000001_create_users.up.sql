-- ==========================================
-- 1.テーブルの作成
-- ==========================================
CREATE TABLE users (
    id UUID PRIMARY KEY,
    user_name TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

-- ==========================================
-- 2.コメント
-- ==========================================
COMMENT ON TABLE users IS 'ユーザー';
COMMENT ON COLUMN users.id IS 'ユーザーID';
COMMENT ON COLUMN users.user_name IS 'ユーザー名';
COMMENT ON COLUMN users.email IS 'メールアドレス';
COMMENT ON COLUMN users.password_hash IS 'ハッシュ化パスワード';
COMMENT ON COLUMN users.created_at IS '作成日時';
COMMENT ON COLUMN users.updated_at IS '更新日時';
COMMENT ON COLUMN users.deleted_at IS '削除日時';

-- ==========================================
-- 3.インデックス・制約
-- ==========================================
-- 論理削除されていないユーザー間でのみ email をユニークにする
CREATE UNIQUE INDEX idx_users_email_active ON users (email) WHERE deleted_at IS NULL;

-- GORMのソフトデリート絞り込み高速化
CREATE INDEX idx_users_deleted_at ON users (deleted_at);
