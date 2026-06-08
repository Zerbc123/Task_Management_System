package model

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	EntityID  uuid.UUID `gorm:"type:uuid"`

	EntityType string
	EventType  string

	Description string

	CreatedAt time.Time
}
