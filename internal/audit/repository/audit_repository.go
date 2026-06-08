package repository

import (
	"context"
	"task-management/internal/audit/model"

	"github.com/google/uuid"
)

type AuditRepository interface {
	Create(ctx context.Context, auditLog *model.AuditLog) error
	GetByEnitityID(ctx context.Context, entityID uuid.UUID) ([]*model.AuditLog, error)
}
