package auth

import (
	"encoding/json"
	"go-todo/gen/api"
	"go-todo/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ハンドラで使うユースケースをモックで定義
type mockAuthUsecase struct {
	signupToken string
	signupUser  *models.User
	signupErr   error
	meUser      *models.User
}

func (m *mockAuthUsecase) Signup(userName, email, password string) (string, *models.User, error) {
	return m.signupToken, m.signupUser, m.signupErr
}

func (m *mockAuthUsecase) Login(email, password string) (string, *models.User, error) {
	return "", nil, nil
}

func (m *mockAuthUsecase) Me(userID uuid.UUID) (*models.User, error) {
	return m.meUser, nil
}

// サインアップ
func TestAuthHandler_Signup_OK(t *testing.T) {
	// given
	// テストモードにして、テスト中のログを抑制する
	gin.SetMode(gin.TestMode)

	// ユースケースをモックで作成
	mockUsecase := &mockAuthUsecase{
		signupToken: "dummy-token",
		signupUser:  &models.User{ID: uuid.New(), Email: "taro@example.com", UserName: "たろう"},
	}
	handler := NewAuthHandler(mockUsecase)

	// レコーダーの準備し、レスポンスステータスやボディを記録する
	recorder := httptest.NewRecorder()
	// recorderに繋がったgin.Contextを作る
	context, _ := gin.CreateTestContext(recorder)

	// リクエストボディの準備
	body := `{"user_name":"たろう","email":"taro@example.com","password":"password123"}`

	// コンテキストに設定
	context.Request = httptest.NewRequest(http.MethodPost, "/api/auth/signup", strings.NewReader(body))
	context.Request.Header.Set("Content-Type", "application/json")

	// when
	// サインアップ処理を実行
	handler.Signup(context)

	// then
	// 201で帰ること
	require.Equal(t, http.StatusCreated, recorder.Code)

	// body を構造体に戻して検証
	var response api.LoginResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.Equal(t, "dummy-token", response.Token)
	assert.Equal(t, "taro@example.com", string(response.User.Email))
}

// サインアップ（エラー）
func TestAuthHandler_Signup_BadRequest(t *testing.T) {
	// given
	// テストモードにして、テスト中のログを抑制する
	gin.SetMode(gin.TestMode)

	// ユースケースをモックで作成
	mockUsecase := &mockAuthUsecase{
		signupToken: "dummy-token",
		signupUser:  &models.User{ID: uuid.New(), Email: "taro@example.com", UserName: "たろう"},
	}
	handler := NewAuthHandler(mockUsecase)

	// レコーダーの準備し、レスポンスステータスやボディを記録する
	recorder := httptest.NewRecorder()
	// recorderに繋がったgin.Contextを作る
	context, _ := gin.CreateTestContext(recorder)

	// リクエストボディの準備
	// わざと変なJSONを作って、バインドをエラーにする
	body := `{"email":"tar`

	// コンテキストに設定
	context.Request = httptest.NewRequest(http.MethodPost, "/api/auth/signup", strings.NewReader(body))
	context.Request.Header.Set("Content-Type", "application/json")

	// when
	// サインアップ処理を実行
	handler.Signup(context)

	// then
	// 400で帰ること
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

// ログアウト
func TestAuthHandler_Logout(t *testing.T) {
	// given
	// テストモードにして、テスト中のログを抑制する
	gin.SetMode(gin.TestMode)

	// ユースケースをモックで作成
	mockUsecase := &mockAuthUsecase{}
	handler := NewAuthHandler(mockUsecase)

	// レコーダーの準備し、レスポンスステータスやボディを記録する
	recorder := httptest.NewRecorder()

	// c.Statusのようなbodyなしレスポンスの場合は、204の書き込みは予約しているだけに過ぎない
	// そのため普通にテストすると200が返ってくる。
	// 上記に対応するために、ginエンジンのServeHTTPを使ってテストする
	_, engine := gin.CreateTestContext(recorder)    // ← エンジンを受け取る
	engine.POST("/api/auth/logout", handler.Logout) // ← ルート登録

	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)

	// when
	// ログアウト処理の実行
	engine.ServeHTTP(recorder, req) // ← エンジン経由で実行(終了時に flush される)

	// then
	// 204で帰ること
	require.Equal(t, http.StatusNoContent, recorder.Code)
}

// GetMe
func TestAuthHandler_GetMe(t *testing.T) {
	// given
	// テストモードにして、テスト中のログを抑制する
	gin.SetMode(gin.TestMode)

	// ユースケースをモックで作成
	mockUsecase := &mockAuthUsecase{
		meUser: &models.User{ID: uuid.New(), Email: "taro@example.com", UserName: "たろう"},
	}
	handler := NewAuthHandler(mockUsecase)

	// レコーダーの準備し、レスポンスステータスやボディを記録する
	recorder := httptest.NewRecorder()
	// recorderに繋がったgin.Contextを作る
	context, _ := gin.CreateTestContext(recorder)

	// コンテキストに設定
	context.Request = httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	// ミドルウェアの代わりにuserIDをここでセット
	context.Set("userID", uuid.New())
	context.Request.Header.Set("Content-Type", "application/json")

	// when
	// GetMe処理の実行
	handler.GetMe(context)

	// then
	// 200で帰ること
	require.Equal(t, http.StatusOK, recorder.Code)

	// ユーザー情報が取得できること
	var response api.User
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.Equal(t, "taro@example.com", string(response.Email))
	assert.Equal(t, "たろう", string(response.UserName))
}
