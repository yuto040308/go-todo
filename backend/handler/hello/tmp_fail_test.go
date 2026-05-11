// このファイルは CI が失敗時にちゃんと red になるかを確認するための一時ファイルです。
// 確認が終わったら必ず削除してください。
package hello

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntentionalFailureForCIVerification(t *testing.T) {
	assert.True(t, false, "CIが失敗時にredになることを確認するための意図的な失敗")
}
