package hello

import (
	"encoding/json"
	"go-todo/handler/testutil"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHelloHandlerOk(t *testing.T) {
	// given
	router := gin.Default()
	router.GET("/hello", HelloHandler)

	// テスト用のリクエストとレコーダーを作成
	request, recorder := testutil.SetupTestRequest("GET", "/hello")

	// when
	router.ServeHTTP(recorder, request)

	// then
	// ステータスコードを検証
	assert.Equal(t, http.StatusOK, recorder.Code)

	// レスポンスボディをJSONとしてパースし、フィールド単位で検証
	var body map[string]string
	err := json.Unmarshal(recorder.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", body["message"])
}
