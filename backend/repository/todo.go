package repository

import (
	"go-todo/models"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *todoRepository {
	return &todoRepository{db: db}
}

// todoRepositoryにCreateメソッドを追加
func (r *todoRepository) Create(todo *models.Todo) error {
	return r.db.Create(todo).Error
}

func (r *todoRepository) FindByUserID(userID uuid.UUID) ([]*models.Todo, error) {
	var todos []*models.Todo
	if err := r.db.Where("user_id = ?", userID).Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *todoRepository) FindByID(id, userID uuid.UUID) (*models.Todo, error) {
	var todo models.Todo
	// 他人のTodoは見れない
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) Update(todo *models.Todo) error {
	if err := r.db.Save(todo).Error; err != nil {
		return err
	}
	return nil
}

func (r *todoRepository) Delete(id, userID uuid.UUID) error {
	var todo models.Todo
	// 他人のTodoは見れない
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		return err
	}
	if err := r.db.Delete(&todo).Error; err != nil {
		return err
	}
	return nil
}
