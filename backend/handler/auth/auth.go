package auth

import (
	"go-todo/gen/api"
	"go-todo/usecase/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type AuthHandler struct {
	authUsecase *auth.AuthUsecase
}

// コンストラクタ
func NewAuthHandler(authUsecase *auth.AuthUsecase) *AuthHandler {
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
