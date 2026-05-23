package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct{
	ID uuid.UUID `json: "id" gorm:"type:uuid;primaryKey"`
	Email string `json: "email" gorm:"unique;not null"`
	PasswordHash string `json:"-"`
	FullName string `json: "full_name"`

	CreatedAt time.Time `json: "created_at"`
	UpdatedAt time.Time `json: "updated_at"`
}