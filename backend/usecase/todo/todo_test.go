package todo

import (
	"go-todo/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ユースケースで使うリポジトリはモックで定義
type mockTodoStore struct {
	listResult []*models.Todo
	// Create に渡された Todo を記録
	createdTodo *models.Todo
}

// Listで使うメソッドをモックにする
func (m *mockTodoStore) FindByUserID(userID uuid.UUID) ([]*models.Todo, error) {
	// 仕込んだ値をそのまま返すだけ
	return m.listResult, nil
}

func (m *mockTodoStore) Create(todo *models.Todo) error {
	m.createdTodo = todo
	return nil
}

// 以下は今回使わないが、インタフェースを満たすために必要(空実装)               { return nil }
func (m *mockTodoStore) FindByID(id, userID uuid.UUID) (*models.Todo, error) { return nil, nil }
func (m *mockTodoStore) Update(todo *models.Todo) error                      { return nil }
func (m *mockTodoStore) Delete(id, userID uuid.UUID) error                   { return nil }

func TestTodoUsecase_List(t *testing.T) {
	// given
	// Listを叩くためのユーザーIDを定義
	userID := uuid.New()

	// モックが2件のTODOを返すように仕込む
	mockTodos := []*models.Todo{
		{Title: "todo1"},
		{Title: "todo2"},
	}
	mockStore := &mockTodoStore{listResult: mockTodos}

	// モックを注入
	usecase := NewTodoUsecase(mockStore)

	// when
	// ユースケースを叩く
	resultTodos, err := usecase.List(userID)

	// then
	// エラーが発生していないこと
	require.NoError(t, err)
	// モックと同じTODOが取得できること
	assert.Equal(t, mockTodos, resultTodos)
}

func TestTodoUsecase_Create(t *testing.T) {
	// given
	userID := uuid.New()

	mockStore := &mockTodoStore{}

	// リポジトリをモックのもので設定しておく
	usecase := NewTodoUsecase(mockStore)

	description := "ネギを買います"

	// when
	resultTodo, err := usecase.Create(userID, "買い物", &description)

	// then
	// エラーが発生してないか
	require.NoError(t, err)

	// 返り値が正しいか
	assert.Equal(t, userID, resultTodo.UserID)
	assert.Equal(t, "買い物", resultTodo.Title)
	assert.Equal(t, description, *resultTodo.Description)

	// store.Create が呼ばれたことだけ確認
	require.NotNil(t, mockStore.createdTodo)
}
