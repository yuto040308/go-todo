package todo

import (
	"encoding/json"
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
type mockTodoUsecase struct {
	getTodo    *models.Todo
	listTodos  []*models.Todo
	createTodo *models.Todo
	updateTodo *models.Todo
}

// モック用のメソッドを定義
func (m *mockTodoUsecase) List(userID uuid.UUID) ([]*models.Todo, error) {
	return m.listTodos, nil
}
func (m *mockTodoUsecase) Create(userID uuid.UUID, title string, description *string) (*models.Todo, error) {
	return m.createTodo, nil
}
func (m *mockTodoUsecase) Delete(id uuid.UUID, userID uuid.UUID) error {
	return nil
}
func (m *mockTodoUsecase) Get(id uuid.UUID, userID uuid.UUID) (*models.Todo, error) {
	return m.getTodo, nil
}
func (m *mockTodoUsecase) Update(id uuid.UUID, userID uuid.UUID, title *string, description *string, isCompleted *bool) (*models.Todo, error) {
	return m.updateTodo, nil
}

// 一覧取得 GET
func TestTodoHandler_List_OK(t *testing.T) {
	// given
	// テストモードにして、テスト中のログを抑制する
	gin.SetMode(gin.TestMode)

	// ユースケースをモックで作成
	mockDescription := "説明"
	mockUsecase := &mockTodoUsecase{
		listTodos: []*models.Todo{
			{
				ID:          uuid.New(),
				UserID:      uuid.New(),
				Title:       "タイトル",
				Description: &mockDescription,
				IsCompleted: true,
			},
		},
	}
	handler := NewTodoHandler(mockUsecase)

	// レコーダーを準備
	recorder := httptest.NewRecorder()
	// recorderに繋がったgin.Contextを作る
	context, _ := gin.CreateTestContext(recorder)
	// コンテキストに設定
	context.Request = httptest.NewRequest(http.MethodGet, "/api/todos", nil)
	// ミドルウェアの代わりにuserIDをここでセット
	context.Set("userID", uuid.New())
	context.Request.Header.Set("Content-Type", "application/json")

	// when
	handler.ListTodos(context)

	// then
	// 200で返ってくること
	require.Equal(t, http.StatusOK, recorder.Code)
	// TODO配列が返ってくること
	var response []*models.Todo
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.Equal(t, "タイトル", response[0].Title)
	assert.Equal(t, "説明", *response[0].Description)
}

// 新規作成 POST
func TestTodoHandler_Create_OK(t *testing.T) {
	// given
	// テストモードにして、テスト中のログを抑制する
	gin.SetMode(gin.TestMode)

	// ユースケースをモックで作成
	mockDescription := "説明"
	mockUsecase := &mockTodoUsecase{
		createTodo: &models.Todo{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			Title:       "タイトル",
			Description: &mockDescription,
			IsCompleted: true,
		},
	}
	handler := NewTodoHandler(mockUsecase)

	// レコーダーを準備
	recorder := httptest.NewRecorder()
	// recorderに繋がったgin.Contextを作る
	context, _ := gin.CreateTestContext(recorder)

	// リクエストボディの準備
	body := `{"title":"タイトル","description":"説明"}`
	// コンテキストに設定
	context.Request = httptest.NewRequest(http.MethodPost, "/api/todos", strings.NewReader(body))
	// ミドルウェアの代わりにuserIDをここでセット
	context.Set("userID", uuid.New())
	context.Request.Header.Set("Content-Type", "application/json")

	// when
	handler.CreateTodo(context)

	// then
	// 201で返ってくること
	require.Equal(t, http.StatusCreated, recorder.Code)
}

// 1件のTODO取得 GET
func TestTodoHandler_Get_OK(t *testing.T) {
	// given
	// テストモードにして、テスト中のログを抑制する
	gin.SetMode(gin.TestMode)

	// ユースケースをモックで作成
	mockDescription := "説明"
	mockUsecase := &mockTodoUsecase{
		getTodo: &models.Todo{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			Title:       "タイトル",
			Description: &mockDescription,
			IsCompleted: true,
		},
	}
	handler := NewTodoHandler(mockUsecase)

	// レコーダーを準備
	recorder := httptest.NewRecorder()
	// recorderに繋がったgin.Contextを作る
	context, _ := gin.CreateTestContext(recorder)
	// コンテキストに設定
	context.Request = httptest.NewRequest(http.MethodGet, "/api/todos/1", nil)
	// ミドルウェアの代わりにuserIDをここでセット
	context.Set("userID", uuid.New())
	context.Request.Header.Set("Content-Type", "application/json")

	// when
	handler.GetTodo(context, uuid.New())

	// then
	// 200で返ってくること
	require.Equal(t, http.StatusOK, recorder.Code)
	// TODOが返ってくること
	var response models.Todo
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.Equal(t, "タイトル", response.Title)
	assert.Equal(t, "説明", *response.Description)
}

// 1件の　TODO追加 Update
func TestTodoHandler_Update_OK(t *testing.T) {
	// given
	// テストモードにして、テスト中のログを抑制する
	gin.SetMode(gin.TestMode)

	// ユースケースをモックで作成
	mockDescription := "説明"
	mockUsecase := &mockTodoUsecase{
		updateTodo: &models.Todo{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			Title:       "タイトル",
			Description: &mockDescription,
			IsCompleted: true,
		},
	}
	handler := NewTodoHandler(mockUsecase)

	// レコーダーを準備
	recorder := httptest.NewRecorder()
	// recorderに繋がったgin.Contextを作る
	context, _ := gin.CreateTestContext(recorder)

	// リクエストボディの準備
	body := `{"title":"タイトル","description":"説明"}`
	// コンテキストに設定
	context.Request = httptest.NewRequest(http.MethodPut, "/api/todos", strings.NewReader(body))
	// ミドルウェアの代わりにuserIDをここでセット
	context.Set("userID", uuid.New())
	context.Request.Header.Set("Content-Type", "application/json")

	// when
	handler.UpdateTodo(context, uuid.New())

	// then
	// 200で返ってくること
	require.Equal(t, http.StatusOK, recorder.Code)
}

// 1件のTODO削除 Delete
func TestTodoHandler_Delete_OK(t *testing.T) {
	// given
	// テストモードにして、テスト中のログを抑制する
	gin.SetMode(gin.TestMode)

	// ユースケースをモックで作成
	mockUsecase := &mockTodoUsecase{}
	handler := NewTodoHandler(mockUsecase)

	// レコーダーを準備
	recorder := httptest.NewRecorder()
	// recorderに繋がったgin.Contextを作る
	_, engine := gin.CreateTestContext(recorder)

	engine.DELETE("/api/todos/1", func(c *gin.Context) {
		c.Set("userID", uuid.New())
		handler.DeleteTodo(c, uuid.New())
	}) // ← ルート登録

	req := httptest.NewRequest(http.MethodDelete, "/api/todos/1", nil)

	// when
	engine.ServeHTTP(recorder, req)

	// then
	// 204を返却
	require.Equal(t, http.StatusNoContent, recorder.Code)
}
