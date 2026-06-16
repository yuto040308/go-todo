package auth

import (
	"go-todo/internal/password"
	"go-todo/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ユースケースで使うリポジトリはモックで定義
type mockUserStore struct {
	// spy: Create に渡された user
	createdUser *models.User
	// stub: FindByEmail が返す user
	findByEmailRes *models.User
	// 見つからないを再現するテスト用
	findByEmailErr error
	// stub: FindByID が返す user
	findByIDRes *models.User
}

func (m *mockUserStore) Create(u *models.User) error {
	m.createdUser = u
	return nil
}
func (m *mockUserStore) FindByEmail(email string) (*models.User, error) {
	return m.findByEmailRes, m.findByEmailErr
}
func (m *mockUserStore) FindByID(id uuid.UUID) (*models.User, error) {
	return m.findByIDRes, nil
}

// サインアップのテスト
func TestAuthUsecase_Signup(t *testing.T) {
	// given
	mockStore := &mockUserStore{}
	usecase := NewAuthUsecase(mockStore)

	// when
	token, user, err := usecase.Signup("ログイン太郎", "taro@example.com", "password1234")

	// then
	require.NoError(t, err)

	// トークンが何かしらの値で発行されること
	assert.NotEmpty(t, token)
	// 指定されたメールアドレスでユーザーが作成されていること
	assert.Equal(t, "taro@example.com", user.Email)
	// ★ 平文ではなくハッシュ化されて保存されているか = Verifyが通れば正しいハッシュ
	require.NotNil(t, mockStore.createdUser)
	// 平文ではない
	assert.NotEqual(t, "password1234", mockStore.createdUser.PasswordHash)
	// 正しいハッシュ
	assert.NoError(t, password.Verify(mockStore.createdUser.PasswordHash, "password1234"))
}

// ログインのテスト
func TestAuthUsecase_Login_成功(t *testing.T) {
	// given
	hashPassword, _ := password.Hash("password1234")
	mockStore := &mockUserStore{
		findByEmailRes: &models.User{Email: "taro@example.com", PasswordHash: hashPassword},
	}
	usecase := NewAuthUsecase(mockStore)

	// when
	token, _, err := usecase.Login("taro@example.com", "password1234")

	// then
	// エラーが起こらない = ログインできており、トークンが発行されていること
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

// ログイン失敗のテスト
func TestAuthUsecase_Login_失敗_パスワード不一致(t *testing.T) {
	// given
	hashPassword, _ := password.Hash("password1234")
	mockStore := &mockUserStore{
		findByEmailRes: &models.User{Email: "taro@example.com", PasswordHash: hashPassword},
	}
	usecase := NewAuthUsecase(mockStore)

	// when
	// わざと間違ったパスワードを指定
	_, _, err := usecase.Login("taro@example.com", "wrong1234")

	// then
	// エラーが発生し、トークンも無効であること
	require.Error(t, err)
}

// ユーザーIDからユーザー情報を取得するテスト
func TestAuthUsecase_Me(t *testing.T) {
	// given
	user := &models.User{Email: "taro@example.com", UserName: "太郎"}
	mockStore := &mockUserStore{findByIDRes: user}
	usecase := NewAuthUsecase(mockStore)

	// when
	resultUser, err := usecase.Me(uuid.New())

	// then
	// エラーが発生してないこと
	require.NoError(t, err)

	// ユーザー情報が取得できること
	assert.Equal(t, user, resultUser)
}
