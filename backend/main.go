package main

import (
	"go-todo/database"
	"go-todo/gen/api"
	authh "go-todo/handler/auth"
	"go-todo/handler/hello"
	todoh "go-todo/handler/todo"
	"go-todo/middleware"
	"go-todo/repository"
	authuc "go-todo/usecase/auth"
	todouc "go-todo/usecase/todo"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/swaggest/swgui/v5emb"
)

type apiServer struct {
	*authh.AuthHandler
	*todoh.TodoHandler
}

func main() {
	// DB接続は handler 実装時 (チケット③) に database.New() を呼び出す
	db, err := database.New()
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}

	// DIのチェーンを手配線する：下から上に
	userRepository := repository.NewUserRepository(db)
	todoRepository := repository.NewTodoRepository(db)

	authUsecase := authuc.NewAuthUsecase(userRepository)
	todoUsecase := todouc.NewTodoUsecase(todoRepository)

	authHandler := authh.NewAuthHandler(authUsecase)
	todoHandler := todoh.NewTodoHandler(todoUsecase)

	// ハンドラをまとめる
	server := &apiServer{AuthHandler: authHandler, TodoHandler: todoHandler}

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

	apiGroup := r.Group("/api")
	apiGroup.GET("/hello", hello.HelloHandler)

	// authとtodo系のルートをapiGroupに登録し、各ルートは認証ミドルウェアを経由するようにまとめて設定
	api.RegisterHandlersWithOptions(apiGroup, server, api.GinServerOptions{
		Middlewares: []api.MiddlewareFunc{
			api.MiddlewareFunc(middleware.AuthHandler()),
		},
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}
