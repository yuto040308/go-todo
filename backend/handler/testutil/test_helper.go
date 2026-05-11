package testutil

import (
	"net/http"
	"net/http/httptest"
)

// setupTestRequest はリクエストとレコーダーを作成します
// method: GET, PUT, DELETE
// url: /todos/:id
func SetupTestRequest(method, url string) (*http.Request, *httptest.ResponseRecorder) {
	request, _ := http.NewRequest(method, url, nil)
	recorder := httptest.NewRecorder()
	return request, recorder
}
