package main

import (
	"go-todo/database"
	"go-todo/handler/hello"
	"go-todo/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 起動時に接続確認
	_ = database.New()

	r := gin.Default()

	r.Use(middleware.CORS())

	api := r.Group("/api")
	api.GET("/hello", hello.HelloHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}
