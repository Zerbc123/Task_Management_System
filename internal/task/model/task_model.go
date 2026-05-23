package model

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusTodo       TaskStatus = "TODO"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusDone       TaskStatus = "DONE"
)

type Task struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	ProjectID   uuid.UUID  `json:"project_id" gorm:"type:uuid;not null"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	AssigneeID  *uuid.UUID `json:"assignee_id,omitempty" gorm:"type:uuid"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}