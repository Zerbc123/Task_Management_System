package repository

import (
	"context"
	"task-management/internal/audit/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type gormAuditRepository struct {
	db *gorm.DB
}

func NewGormAuditRepository(db *gorm.DB) AuditRepository {
	return &gormAuditRepository{db: db}
}

func (r *gormAuditRepository) Create(ctx context.Context, auditLog *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(auditLog).Error
}

func (r *gormAuditRepository) GetByEnitityID(ctx context.Context, entityID uuid.UUID) ([]*model.AuditLog, error) {
	var logs []*model.AuditLog

	err := r.db.WithContext(ctx).Where("entity_id = ?", entityID).Order("created_at asc").Find(&logs).Error

	return logs, err
}
