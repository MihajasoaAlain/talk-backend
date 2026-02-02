package repository

import (
	"talk-backend/internal/models"

	"gorm.io/gorm"
)

type AuditRepository interface {
	Create(log *models.AuditLog) error
}

type auditRepository struct{ db *gorm.DB }

func NewAuditRepository(db *gorm.DB) AuditRepository { return &auditRepository{db: db} }

func (r *auditRepository) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}
