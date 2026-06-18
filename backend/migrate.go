package main

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	// postgres:// スキームの DB ドライバを登録 (blank import で副作用のみ)
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// migrationFiles はマイグレーション SQL を Go バイナリに埋め込む。
// Cloud Run はローカルファイルをマウントできないため、ファイルパス依存をなくす。
//
//go:embed migrations/*.sql
var migrationFiles embed.FS

// runMigrations は起動時に「テーブルが無ければ作る」= migrate up を実行する。
// 適用済みなら ErrNoChange なので no-op (= 何もしない)。
func runMigrations(dsn string) error {
	// 埋め込んだ SQL を migrate のソースとして読む
	src, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("マイグレーションソースの読み込み失敗: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, dsn)
	if err != nil {
		return fmt.Errorf("migrate の初期化失敗: %w", err)
	}

	// up を適用。既に最新なら ErrNoChange が返るので、それは成功扱いにする
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up 失敗: %w", err)
	}
	return nil
}
