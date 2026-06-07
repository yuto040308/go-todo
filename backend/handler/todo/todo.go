package todo

import (
	"errors"
	"go-todo/gen/api"
	"go-todo/models"
	"go-todo/usecase/todo"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TODO: 5メソッド(ListTodos / CreateTodo / GetTodo / UpdateTodo / DeleteTodo)

// DTO変換ヘルパ
func toTodoDTO(t *models.Todo) api.Todo {
	return api.Todo{
		Id:          t.ID,
		UserId:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		IsCompleted: t.IsCompleted,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// 構造体
type TodoHandler struct {
	todoUsecase *todo.TodoUsecase
}

// コンストラクタ
func NewTodoHandler(todoUsecase *todo.TodoUsecase) *TodoHandler {
	return &TodoHandler{todoUsecase: todoUsecase}
}

// GET /todos 一覧取得
func (h *TodoHandler) ListTodos(c *gin.Context) {
	// ミドルウェアからuserIDを取得
	userID := c.MustGet("userID").(uuid.UUID)

	// userIDをキーに紐づく全てのtodoを取得
	todos, err := h.todoUsecase.List(userID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			api.Error{Code: "INTERNAL_SERVER_ERROR", Message: "Todoの取得に失敗しました"},
		)
		return
	}

	// 変換して返却
	// makeにして0件でも空配列が帰るようにする
	convertTodos := make([]api.Todo, 0, len(todos))
	for _, todo := range todos {
		convertTodos = append(convertTodos, toTodoDTO(todo))
	}

	c.JSON(http.StatusOK, convertTodos)
}

// POST /todos 新規作成
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	// ミドルウェアからuserIDを取得
	userID := c.MustGet("userID").(uuid.UUID)

	// リクエストをバインド
	var req api.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.Error{Code: "VALIDATION_ERROR", Message: err.Error()})
		return
	}

	// 保存
	todo, err := h.todoUsecase.Create(userID, req.Title, req.Description)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			api.Error{Code: "INTERNAL_SERVER_ERROR", Message: "Todoの作成に失敗しました"},
		)
		return
	}

	// OKのレスポンス
	c.JSON(http.StatusCreated, toTodoDTO(todo))
}

// GET /todos/{id}
// router.GET(options.BaseURL+"/todos/:id", wrapper.GetTodo) を経由するので、idが取れる
func (h *TodoHandler) GetTodo(c *gin.Context, todoID api.TodoId) {
	// ミドルウェアからuserIDを取得
	userID := c.MustGet("userID").(uuid.UUID)

	// 取得
	todo, err := h.todoUsecase.Get(todoID, userID)
	if err != nil {
		// エラーコードにより分岐
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound,
				api.Error{Code: "NOT_FOUND", Message: "Todoが見つかりませんでした"},
			)
		} else {
			c.JSON(
				http.StatusInternalServerError,
				api.Error{Code: "INTERNAL_SERVER_ERROR", Message: "Todoの取得に失敗しました"},
			)
		}
		return
	}

	// 変換して返却
	c.JSON(http.StatusOK, toTodoDTO(todo))
}

// PUT /todos/{id}

// DELETE /todos/{id}
