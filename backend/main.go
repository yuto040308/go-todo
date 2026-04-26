package main

import (
	"go-todo/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(middleware.CORS())

	api := r.Group("/api")
	api.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!!!!!!!!!"})
	})

	r.Run(":8080")
}
