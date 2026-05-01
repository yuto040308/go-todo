package main

import "os"

// CIの動作確認用のサンプル。動作が確認できたらこのファイルは削除してください。
//
// このファイルでは下記の2種類の指摘が同時に発生します:
//   - 余分な空行 → gofumpt が自動修正 → reviewdog が Commit suggestion として提示
//   - os.Setenv の戻り値(error)を握りつぶし → errcheck が検知 → 最後のステップでCIが失敗
func RunLintDemo() {
	os.Setenv("DEMO_KEY", "demo_value")
}
