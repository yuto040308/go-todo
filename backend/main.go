package main

import (
	"go-todo/handler/hello"
	"go-todo/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/swaggest/swgui/v5emb"
)

func main() {
	// DB接続は handler 実装時 (チケット③) に database.New() を呼び出す

	r := gin.Default()

	// 開発環境のみ Swagger UIを配信（本番では切っておく想定）
	if os.Getenv("APP_ENV") != "production" {
		// openapi.yamlでOpenAPI仕様ファイル自体を配信 これがないとSwagger UIが動かない
		r.StaticFile("/openapi.yaml", "/api/openapi.yaml")

		// /swagger/ 配下でSwagger UIを配信
		r.GET("/swagger/*any", gin.WrapH(
			v5emb.New("go-todo API", "/openapi.yaml", "/swagger/"),
		))
	}

	r.Use(middleware.CORS())

	api := r.Group("/api")
	api.GET("/hello", hello.HelloHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}
