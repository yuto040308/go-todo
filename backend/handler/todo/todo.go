package todo

import (
	"errors"
	"go-todo/gen/api"
	"go-todo/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 5メソッド(ListTodos / CreateTodo / GetTodo / UpdateTodo / DeleteTodo)

// TodoUsecase は handler が必要とする usecase 操作 (consumer 側で定義 = テストでモック可能)
type TodoUsecase interface {
	List(userID uuid.UUID) ([]*models.Todo, error)
	Create(userID uuid.UUID, title string, description *string) (*models.Todo, error)
	Get(id, userID uuid.UUID) (*models.Todo, error)
	Update(id, userID uuid.UUID, title, description *string, isCompleted *bool) (*models.Todo, error)
	Delete(id, userID uuid.UUID) error
}

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
	todoUsecase TodoUsecase
}

// コンストラクタ
func NewTodoHandler(todoUsecase TodoUsecase) *TodoHandler {
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
func (h *TodoHandler) UpdateTodo(c *gin.Context, todoID api.TodoId) {
	// ミドルウェアからuserIDを取得
	userID := c.MustGet("userID").(uuid.UUID)

	// リクエストパラメーターをバインド
	var req api.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.Error{Code: "VALIDATION_ERROR", Message: err.Error()})
		return
	}

	// ユースケースを使って更新
	todo, err := h.todoUsecase.Update(
		todoID,
		userID,
		req.Title,
		req.Description,
		req.IsCompleted,
	)
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
				api.Error{Code: "INTERNAL_SERVER_ERROR", Message: "Todoの更新に失敗しました"},
			)
		}
		return
	}

	// 更新後のTodoを返す
	c.JSON(http.StatusOK, toTodoDTO(todo))
}

// DELETE /todos/{id}
func (h *TodoHandler) DeleteTodo(c *gin.Context, todoID api.TodoId) {
	// ミドルウェアからuserIDを取得
	userID := c.MustGet("userID").(uuid.UUID)

	// 削除
	if err := h.todoUsecase.Delete(todoID, userID); err != nil {
		// エラーコードにより分岐
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound,
				api.Error{Code: "NOT_FOUND", Message: "Todoが見つかりませんでした"},
			)
		} else {
			c.JSON(
				http.StatusInternalServerError,
				api.Error{Code: "INTERNAL_SERVER_ERROR", Message: "Todoの削除に失敗しました"},
			)
		}
		return
	}

	// 204を返却
	c.Status(http.StatusNoContent)
}
