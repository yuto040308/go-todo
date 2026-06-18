package testutil

import (
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// テストで使うヘルパを定義

// テストDBに接続するヘルパ
func SetupTestDB(t *testing.T) *gorm.DB {
	// 失敗時に行番号を呼び出し元に戻す
	// これが書いてあると、エラーの時にこのメソッドではなく呼び出し元がエラー判定として表示されるので便利
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		// 明示的に落とす
		t.Fatal("TEST_DATABASE_URL が未設定です")
	}

	// DBに接続
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// テスト中の SQL/エラーログを抑制
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("テストDB接続失敗: %v", err)
	}

	return db
}

// テーブルを空にするヘルパ。usersもtodosも両方空にする
func TruncateAll(t *testing.T, db *gorm.DB) {
	t.Helper()

	// CASCADEでFK依存も一緒に消す
	// RESTART IDENTITY は連番リセット
	if err := db.Exec("TRUNCATE TABLE todos, users RESTART IDENTITY CASCADE").Error; err != nil {
		t.Fatalf("truncate 失敗: %v", err)
	}
}
