-- ==========================================
-- 1.テーブルの作成
-- ==========================================
CREATE TABLE todos (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,

    -- 外部キー制約（TODOがあるユーザーが削除されるのを防ぐ）
    CONSTRAINT fk_todos_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
);

-- ==========================================
-- 2.コメント
-- ==========================================
COMMENT ON TABLE todos IS 'TODO';
COMMENT ON COLUMN todos.id IS 'TODO ID';
COMMENT ON COLUMN todos.user_id IS 'ユーザーID';
COMMENT ON COLUMN todos.title IS 'タイトル';
COMMENT ON COLUMN todos.description IS '説明';
COMMENT ON COLUMN todos.is_completed IS '完了フラグ';
COMMENT ON COLUMN todos.created_at IS '作成日時';
COMMENT ON COLUMN todos.updated_at IS '更新日時';
COMMENT ON COLUMN todos.deleted_at IS '削除日時';

-- ==========================================
-- 3.インデックス・制約
-- ==========================================
-- 外部キーの高速化
CREATE INDEX idx_todos_user_id ON todos (user_id);

-- GORMのソフトデリート絞り込み高速化
CREATE INDEX idx_todos_deleted_at ON todos (deleted_at);