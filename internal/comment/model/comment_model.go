package model

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	TaskID  uuid.UUID `json:"task_id" gorm:"type:uuid;not null"`
	UserID  uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Content string    `json:"content" gorm:"not null"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
