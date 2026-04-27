package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 通常はNginx経由(http://localhost)でアクセスするためCORSは発生しない。
		// Nginxを介さず直接バックエンドを叩くケースのために許可オリジンを定義しておく。
		AllowOrigins:     []string{"http://localhost", "http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	})
}
