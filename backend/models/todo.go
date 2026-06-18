package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Todo は GORM の DB モデル。User とは belongs-to の関係 (UserID は外部キー)。
type Todo struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	Title       string         `json:"title" gorm:"not null" binding:"required"`
	Description *string        `json:"description"`
	IsCompleted bool           `json:"is_completed" gorm:"column:is_completed;default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (t *Todo) BeforeCreate(_ *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
