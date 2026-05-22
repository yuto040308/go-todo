package main

import (
	"go-todo/database"
	"go-todo/handler/hello"
	"go-todo/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 起動時にDB接続確認 (失敗時は起動を中断)
	db, err := database.New()
	if err != nil {
		log.Fatalf("データベース接続に失敗: %v", err)
	}
	_ = db // 後の CRUD チケットで handler に渡す

	r := gin.Default()

	r.Use(middleware.CORS())

	api := r.Group("/api")
	api.GET("/hello", hello.HelloHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}
