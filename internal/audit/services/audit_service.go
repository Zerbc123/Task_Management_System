package services

import (
	"context"
	"task-management/internal/audit/model"
	"task-management/internal/audit/repository"
	"time"

	"github.com/google/uuid"
)

type AuditService interface {
	Log(ctx context.Context, userID uuid.UUID, entityID uuid.UUID, entityType string, eventType string, description string) error
	GetByEntityID(ctx context.Context, entityID uuid.UUID) ([]*model.AuditLog, error)
}

type auditService struct {
	repo repository.AuditRepository
}

func NewAuditService(repo repository.AuditRepository) AuditService {
	return &auditService{repo: repo}
}

func (s *auditService) Log(ctx context.Context, userID uuid.UUID, entityID uuid.UUID, entityType string, eventType string, description string) error {
	auditLog := &model.AuditLog{
		ID: uuid.New(),
		UserID: userID,
		EntityID: entityID,
		EntityType: entityType,
		EventType: eventType,
		Description: description,
		CreatedAt: time.Now(),
	}

	return s.repo.Create(ctx, auditLog)
}

func (s *auditService) GetByEntityID(ctx context.Context, entityID uuid.UUID) ([]*model.AuditLog, error) {
	return s.repo.GetByEnitityID(ctx, entityID)
}
