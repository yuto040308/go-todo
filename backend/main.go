package main

import (
	"context"
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

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	oapimw "github.com/oapi-codegen/gin-middleware"
	"github.com/swaggest/swgui/v5emb"
)

type apiServer struct {
	*authh.AuthHandler
	*todoh.TodoHandler
}

func main() {
	// 起動時に migrate up (テーブルが無ければ作る。適用済みなら no-op)。
	// 失敗しても Cloud Run のブルーグリーンで旧リビジョンが残るため本番は停止しない。
	if err := runMigrations(os.Getenv("DATABASE_URL")); err != nil {
		log.Fatalf("マイグレーション失敗: %v", err)
	}

	// DB接続
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

	// OpenAPI spec を読み込み、リクエストを spec と照合する validator を作る
	swagger, err := api.GetSpec()
	if err != nil {
		log.Fatalf("OpenAPI spec の読み込み失敗: %v", err)
	}
	// host に依存せず basePath /api だけで照合させる (server URL を相対パスにする)
	swagger.Servers = openapi3.Servers{{URL: "/api"}}

	// validator はリクエストを OpenAPI spec と照合する門番。呼ばれる流れ:
	//   リクエスト到着
	//     ├─ security 必須ルート? → AuthenticationFunc を呼ぶ →(nil)→ 通す(認証は自前 AuthHandler 担当)
	//     ├─ body が spec 違反?   → ErrorHandler を呼ぶ → 400 {code,message} で終了
	//     └─ 全部 OK              → 次の middleware(auth gate)へ
	validator := oapimw.OapiRequestValidatorWithOptions(swagger, &oapimw.Options{
		Options: openapi3filter.Options{
			// security(bearerAuth)の検証は自前の AuthHandler に任せるのでここでは素通し
			AuthenticationFunc: func(_ context.Context, _ *openapi3filter.AuthenticationInput) error {
				return nil
			},
		},
		// バリデーターがスペック違反を判断した場合のレスポンスの書き方を定義
		// エラー形式をプロジェクト規約 {code, message} に揃える
		ErrorHandler: func(c *gin.Context, message string, statusCode int) {
			c.JSON(statusCode, api.Error{Code: "VALIDATION_ERROR", Message: message})
		},
	})

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
			// ① まず spec と照合 (必須/形式バリデーション)
			api.MiddlewareFunc(validator),
			// ② 次に認証 (BearerAuthScopes が立ってるルートだけ)
			api.MiddlewareFunc(middleware.AuthHandler()),
		},
	})

	// Cloud Run は $PORT でリッスンするポートを指定してくる。ローカルは 8080。
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}
