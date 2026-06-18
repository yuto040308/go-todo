package middleware

import (
	"go-todo/gen/api"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthHandler は Authorization: Bearer <JWT> を検証する gin ミドルウェア。
// JWT の claim から user_id (UUID 文字列) を取り出し、c.Set("userID", uuid.UUID) で
// 後続ハンドラに渡す。エラー時は {code, message} 形式の 401 で abort する。
func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// BearerAuthScopesに目印がついてない = 認証不要の場合はそのまま通す
		if _, required := c.Get(string(api.BearerAuthScopes)); !required {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			abortUnauthorized(c, "Authorization header is missing")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			abortUnauthorized(c, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			abortUnauthorized(c, "Invalid token claims")
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			abortUnauthorized(c, "user_id claim missing")
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			abortUnauthorized(c, "Invalid user_id format")
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func abortUnauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"code":    "UNAUTHORIZED",
		"message": message,
	})
}
