package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User は GORM の DB モデル。API DTO とは別物 (API DTO は gen/api/ で生成される)。
// PasswordHash は json:"-" でレスポンスから完全に隠す。
type User struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	UserName     string         `json:"user_name" gorm:"not null"`
	Email        string         `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string         `json:"-" gorm:"not null"`
	Todos        []Todo         `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate は GORM の hook。Create 直前に UUID を自動採番する。
func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
