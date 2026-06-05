package auth

import (
	"go-todo/internal/password"
	"go-todo/internal/token"
	"go-todo/models"

	"github.com/google/uuid"
)

// 1.ユースケースが必要なインターフェースを定義
type UserStore interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uuid.UUID) (*models.User, error)
}

// 2.ユースケース本体
type AuthUsecase struct {
	userStore UserStore
}

// 3.コンストラクタ
func NewAuthUsecase(userStore UserStore) *AuthUsecase {
	return &AuthUsecase{userStore: userStore}
}

// 4.サインアップ（新規登録）
func (u *AuthUsecase) Signup(userName, email, plainPassword string) (string, *models.User, error) {
	// 1. password.Hash(plainPassword) でハッシュ化
	hashedPassword, err := password.Hash(plainPassword)
	if err != nil {
		return "", nil, err
	}

	// 2. &models.User{UserName: ..., Email: ..., PasswordHash: ハッシュ} を組み立て
	user := models.User{
		UserName:     userName,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	// 3. u.userStore.Create(user) で保存(ID は BeforeCreate が採番)
	if err := u.userStore.Create(&user); err != nil {
		return "", nil, err
	}

	// 4. token.Generate(user.ID) で JWT 発行
	jwtToken, err := token.Generate(user.ID)
	if err != nil {
		return "", nil, err
	}

	// 5. return token, user, nil
	return jwtToken, &user, nil
}
