package repository

import (
	"talk-backend/internal/models"

	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(msg *models.Message) error
	List(conversationID uint, limit int, beforeID *uint) ([]models.Message, error)
}

type messageRepository struct{ db *gorm.DB }

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(msg *models.Message) error {
	return r.db.Create(msg).Error
}

func (r *messageRepository) List(conversationID uint, limit int, beforeID *uint) ([]models.Message, error) {
	if limit <= 0 || limit > 100 {
		limit = 30
	}

	q := r.db.Where("conversation_id = ?", conversationID).Order("id DESC").Limit(limit)
	if beforeID != nil && *beforeID > 0 {
		q = q.Where("id < ?", *beforeID)
	}

	var msgs []models.Message
	err := q.Find(&msgs).Error
	return msgs, err
}
