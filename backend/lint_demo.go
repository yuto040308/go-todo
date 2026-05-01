package main

import (
	"fmt"
	"os"
)

// CIの動作確認用のサンプル。動作が確認できたらこのファイルは削除してください。
//
// このファイルには errcheck が指摘する「エラー戻り値の握りつぶし」だけが含まれています。
// errcheck は自動修正できない種類なので reviewdog はボタンなしのインラインコメントで指摘し、
// 最終ステップでCIが赤になります。
// 自分でエラーハンドリングを書いて push し直すとCIが緑になります。
func RunLintDemo() {
	os.Setenv("DEMO_KEY", "demo_value")
	fmt.Println("CI lint demo")
}
