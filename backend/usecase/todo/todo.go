package todo

import (
	"go-todo/models"

	"github.com/google/uuid"
)

// 1.ユースケースが必要なインターフェースを定義
type TodoStore interface {
	Create(todo *models.Todo) error
	FindByUserID(userID uuid.UUID) ([]*models.Todo, error)
	FindByID(id, userID uuid.UUID) (*models.Todo, error)
	Update(todo *models.Todo) error
	Delete(id, userID uuid.UUID) error
}

// 2.ユースケース本体
type TodoUsecase struct {
	todoStore TodoStore
}

// 3.コンストラクタ
func NewTodoUsecase(todoStore TodoStore) *TodoUsecase {
	return &TodoUsecase{todoStore: todoStore}
}

// 4.Create
// *はポインタではなく、任意という意味で使う
func (u *TodoUsecase) Create(userID uuid.UUID, title string, description *string) (*models.Todo, error) {
	// TODOに詰め替える
	todo := models.Todo{
		UserID:      userID,
		Title:       title,
		Description: description,
	}

	// 保存 -> 参照渡しで他のカラムの値を埋めてもらう
	if err := u.todoStore.Create(&todo); err != nil {
		return nil, err
	}

	// 作成されたTodoを返却
	return &todo, nil
}

// 5.List
func (u *TodoUsecase) List(userID uuid.UUID) ([]*models.Todo, error) {
	// ユーザーIDをキーにして紐づくTodoを全件取得
	todos, err := u.todoStore.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 返却
	return todos, nil
}

// 6.Get
func (u *TodoUsecase) Get(id, userID uuid.UUID) (*models.Todo, error) {
	// IDとユーザーIDをキーにして紐づくTodoを取得
	todo, err := u.todoStore.FindByID(id, userID)
	if err != nil {
		return nil, err
	}

	// 返却
	return todo, nil
}

// 7.Update
func (u *TodoUsecase) Update(id, userID uuid.UUID, title, description *string, isCompleted *bool) (*models.Todo, error) {
	// IDとユーザーIDをキーにして紐づくTodoを取得
	todo, err := u.todoStore.FindByID(id, userID)
	if err != nil {
		return nil, err
	}

	// 値を入れ替え
	// 任意なので値があればセットする
	if title != nil {
		todo.Title = *title
	}
	if description != nil {
		todo.Description = description
	}
	if isCompleted != nil {
		todo.IsCompleted = *isCompleted
	}

	// 保存
	if err := u.todoStore.Update(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

// 8.Delete
func (u *TodoUsecase) Delete(id, userID uuid.UUID) error {
	if err := u.todoStore.Delete(id, userID); err != nil {
		return err
	}
	return nil
}
