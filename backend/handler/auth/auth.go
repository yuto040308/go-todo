package auth

import (
	"go-todo/gen/api"
	"go-todo/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type AuthUsecase interface {
	Signup(userName, email, password string) (string, *models.User, error)
	Login(email, password string) (string, *models.User, error)
	Me(userID uuid.UUID) (*models.User, error)
}

type AuthHandler struct {
	authUsecase AuthUsecase
}

// コンストラクタ
func NewAuthHandler(authUsecase AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

// サインアップ関数
func (h *AuthHandler) Signup(c *gin.Context) {
	// ① リクエストをバインド
	var req api.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.Error{Code: "VALIDATION_ERROR", Message: err.Error()})
		return
	}

	// ② usecase 呼び出し(req.Email は openapi_types.Email 型なので string() 変換)
	token, user, err := h.authUsecase.Signup(req.UserName, string(req.Email), req.Password)
	if err != nil {
		// email重複等
		c.JSON(http.StatusConflict, api.Error{Code: "CONFLICT", Message: err.Error()})
		return
	}

	// ③④ DTO 変換してレスポンス(201)
	c.JSON(http.StatusCreated, api.LoginResponse{
		Token: token,
		User: api.User{
			Id:       user.ID,
			Email:    openapi_types.Email(user.Email),
			UserName: user.UserName,
		},
	})
}

// ログイン
func (h *AuthHandler) Login(c *gin.Context) {
	// リクエストをバインド
	var req api.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.Error{Code: "VALIDATION_ERROR", Message: err.Error()})
		return
	}

	// ログイン作業
	token, user, err := h.authUsecase.Login(string(req.Email), req.Password)
	if err != nil {
		// 「メール無い」「パスワード違う」を区別せず一律 401
		c.JSON(http.StatusUnauthorized, api.Error{Code: "UNAUTHORIZED", Message: "email or password"})
		return
	}

	// レスポンス作成
	c.JSON(http.StatusOK, api.LoginResponse{
		Token: token,
		User: api.User{
			Id:       user.ID,
			Email:    openapi_types.Email(user.Email),
			UserName: user.UserName,
		},
	})
}

// ログアウト
// ステートレスなのでサーバー側は204を返すだけ
func (h *AuthHandler) Logout(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// context の userID からユーザー情報を取得
func (h *AuthHandler) GetMe(c *gin.Context) {
	// ミドルウェアが入れた値を取得
	userID := c.MustGet("userID").(uuid.UUID)

	// ユーザー情報を取得
	user, err := h.authUsecase.Me(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, api.Error{Code: "NOT_FOUND", Message: "user not found"})
		return
	}

	c.JSON(http.StatusOK, api.User{
		Id:       user.ID,
		Email:    openapi_types.Email(user.Email),
		UserName: user.UserName,
	})
}
