package repository

import (
	"go-todo/models"
	"go-todo/testutil"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// Create
func TestTodoRepository_Create(t *testing.T) {
	// given
	// テストのDB接続 + まっさらにする
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db)

	// リポジトリの作成
	repository := NewTodoRepository(db)

	// ユーザーを保存しておく
	user := &models.User{UserName: "たろう", Email: "taro@example.com", PasswordHash: "hashed"}
	require.NoError(t, db.Create(user).Error)

	// TODOを用意
	description := "説明"
	todo := &models.Todo{UserID: user.ID, Title: "タイトル", Description: &description}

	// when
	// TODOをDBに保存
	err := repository.Create(todo)

	// then
	// エラーが発生していないこと
	require.NoError(t, err)

	// 中身が問題ないこと
	require.NotEqual(t, uuid.Nil, todo.ID)
	assert.Equal(t, user.ID, todo.UserID)
	assert.Equal(t, "タイトル", todo.Title)
	assert.Equal(t, "説明", *todo.Description)
	assert.Equal(t, false, todo.IsCompleted) // デフォルトはfalse
}

// FindByUserID（自分のTODOを取得）
func TestTodoRepository_FindByUserID(t *testing.T) {
	// given
	// テストのDB接続 + まっさらにする
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db)

	// リポジトリを作成
	repository := NewTodoRepository(db)

	// ユーザーを保存しておく
	user := &models.User{UserName: "たろう", Email: "taro@example.com", PasswordHash: "hashed"}
	require.NoError(t, db.Create(user).Error)

	// TODOも保存しておく
	description := "説明"
	todo := &models.Todo{UserID: user.ID, Title: "タイトル", Description: &description}
	require.NoError(t, db.Create(todo).Error)

	// when
	// ユーザーをキーにして取得する
	resultTodos, err := repository.FindByUserID(user.ID)

	// then
	require.NoError(t, err)
	// 正しく1件取得できる
	assert.Len(t, resultTodos, 1)
	assert.NotEqual(t, uuid.Nil, resultTodos[0].ID)
	assert.Equal(t, user.ID, resultTodos[0].UserID)
	assert.Equal(t, todo.Title, resultTodos[0].Title)
}

// FindByID
func TestTodoRepository_FindByID(t *testing.T) {
	// given
	// テストのDB接続 + まっさらにする
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db)

	// リポジトリを作成
	repository := NewTodoRepository(db)

	// ユーザーを保存しておく
	user := &models.User{UserName: "たろう", Email: "taro@example.com", PasswordHash: "hashed"}
	require.NoError(t, db.Create(user).Error)

	// TODOも保存しておく
	description := "説明"
	todo := &models.Todo{UserID: user.ID, Title: "タイトル", Description: &description}
	require.NoError(t, db.Create(todo).Error)

	// when
	// ユーザーをキーにして取得する
	resultTodo, err := repository.FindByID(todo.ID, user.ID)

	// then
	require.NoError(t, err)
	// 正しくTODOが取得できる
	assert.NotEqual(t, uuid.Nil, resultTodo.ID)
	assert.Equal(t, user.ID, resultTodo.UserID)
	assert.Equal(t, todo.Title, resultTodo.Title)
}

// FindByID(他人のものは取れない)
func TestTodoRepository_FindByID_他人のものは取得できない(t *testing.T) {
	// given
	// テストのDB接続 + まっさらにする
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db)

	// リポジトリを作成
	repository := NewTodoRepository(db)

	// ユーザーAとそのTODOを用意
	userA := &models.User{UserName: "たろう", Email: "taro@example.com", PasswordHash: "hashed"}
	require.NoError(t, db.Create(userA).Error)
	description := "説明A"
	todoA := &models.Todo{UserID: userA.ID, Title: "タイトルA", Description: &description}
	require.NoError(t, db.Create(todoA).Error)

	// ユーザーBを用意
	userB := &models.User{UserName: "じろう", Email: "jiro@example.com", PasswordHash: "hashed"}
	require.NoError(t, db.Create(userB).Error)

	// when
	// ユーザーBがユーザーAのTODOを取ろうとする
	_, err := repository.FindByID(todoA.ID, userB.ID)

	// then
	// エラーで取れない
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

// Update
func TestTodoRepository_Update(t *testing.T) {
	// given
	// テストのDB接続 + まっさらにする
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db)

	// リポジトリを作成
	repository := NewTodoRepository(db)

	// ユーザーを保存しておく
	user := &models.User{UserName: "たろう", Email: "taro@example.com", PasswordHash: "hashed"}
	require.NoError(t, db.Create(user).Error)

	// TODOも保存しておく
	description := "説明"
	todo := &models.Todo{UserID: user.ID, Title: "タイトル", Description: &description}
	require.NoError(t, db.Create(todo).Error)

	// when
	// 更新する
	todo.Title = "変更後タイトル"
	updateDescription := "変更後説明"
	todo.Description = &updateDescription
	err := repository.Update(todo)

	// then
	require.NoError(t, err)

	// DBから再度取り直して、変わっていることを確認
	updated, err := repository.FindByID(todo.ID, user.ID) // ← 新しく SELECT し直す
	require.NoError(t, err)
	assert.Equal(t, "変更後タイトル", updated.Title) // DB の値が新しい
	assert.Equal(t, "変更後説明", *updated.Description)
}

// Delete
func TestTodoRepository_Delete(t *testing.T) {
	// given
	// テストのDB接続 + まっさらにする
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db)

	// リポジトリを作成
	repository := NewTodoRepository(db)

	// ユーザーを保存しておく
	user := &models.User{UserName: "たろう", Email: "taro@example.com", PasswordHash: "hashed"}
	require.NoError(t, db.Create(user).Error)

	// TODOも保存しておく
	description := "説明"
	todo := &models.Todo{UserID: user.ID, Title: "タイトル", Description: &description}
	require.NoError(t, db.Create(todo).Error)

	// when
	// 削除する
	err := repository.Delete(todo.ID, user.ID)

	// then
	require.NoError(t, err)
	// TODOが消えていること（ソフトデリートされて取得できない）
	_, findErr := repository.FindByID(todo.ID, user.ID)
	require.ErrorIs(t, findErr, gorm.ErrRecordNotFound)
}
