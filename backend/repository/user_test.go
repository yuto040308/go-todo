package repository

import (
	"go-todo/models"
	"go-todo/testutil"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 新規作成
func TestUserRepository_Create(t *testing.T) {
	// given
	// テストのDB接続 + まっさらにする
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db)

	// リポジトリを作成
	repository := NewUserRepository(db)

	// ユーザーを作成
	user := &models.User{UserName: "たろう", Email: "taro@example.com", PasswordHash: "hashed"}

	// when
	// ユーザーをDBに保存
	err := repository.Create(user)

	// then
	require.NoError(t, err)
	// uuidがセットされており、値もセットされている
	require.NotEqual(t, uuid.Nil, user.ID)
	assert.Equal(t, "たろう", user.UserName)
	assert.Equal(t, "taro@example.com", user.Email)
	assert.False(t, user.CreatedAt.IsZero())
}

// email重複（同じ email で Create 2回 → 2回目がエラー。部分ユニーク制約が効く証明)
func TestUserRepository_Create_Email重複(t *testing.T) {
	// given
	// テストのDB接続 + まっさらにする
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db)

	// リポジトリを作成
	repository := NewUserRepository(db)

	// ユーザーを作成
	user1 := &models.User{UserName: "たろう", Email: "taro@example.com", PasswordHash: "hashed"}

	// when1
	// ユーザーをDBに保存
	err := repository.Create(user1)

	// then1
	require.NoError(t, err)

	// when2
	// 同じメールアドレスで再度保存
	user2 := &models.User{UserName: "じろう", Email: "taro@example.com", PasswordHash: "hashed"}
	errTwo := repository.Create(user2)

	// 同じメールアドレスのため保存できない
	require.Error(t, errTwo)
}

// FindByEmail
func TestUserRepository_FindByEmail(t *testing.T) {
	// given
	// テストのDB接続 + まっさらにする
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db)

	// リポジトリを作成
	repository := NewUserRepository(db)

	// ユーザーを作成
	user := &models.User{UserName: "たろう", Email: "taro@example.com", PasswordHash: "hashed"}

	// 手動でDBに一度保存する
	require.NoError(t, db.Create(user).Error)

	// when
	// ユーザーをEmailをキーにして取得する
	resultUser, err := repository.FindByEmail("taro@example.com")

	// then
	require.NoError(t, err)
	// 正しく取得できる
	assert.Equal(t, "たろう", resultUser.UserName)
	assert.Equal(t, "taro@example.com", resultUser.Email)
}

// FindByID
func TestUserRepository_FindByID(t *testing.T) {
	// given
	// テストのDB接続 + まっさらにする
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db)

	// リポジトリを作成
	repository := NewUserRepository(db)

	// ユーザーを作成
	user := &models.User{UserName: "たろう", Email: "taro@example.com", PasswordHash: "hashed"}

	// 手動でDBに一度保存する
	require.NoError(t, db.Create(user).Error)

	// when
	// ユーザーをEmailをキーにして取得する
	resultUser, err := repository.FindByID(user.ID)

	// then
	require.NoError(t, err)
	// 正しく取得できる
	assert.Equal(t, "たろう", resultUser.UserName)
	assert.Equal(t, "taro@example.com", resultUser.Email)
}
