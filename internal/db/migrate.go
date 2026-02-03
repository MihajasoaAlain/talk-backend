package db

import (
	"talk-backend/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.RefreshToken{}, &models.AuditLog{}, &models.Conversation{}, &models.ConversationMember{}, &models.Message{})
}
