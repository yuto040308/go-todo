package main

import (
	"go-todo/handler/hello"
	"go-todo/middleware"
	"log"

	"go-todo/handler/hello"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(middleware.CORS())

	api := r.Group("/api")
	api.GET("/hello", hello.HelloHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}
