package repository

import (
	"talk-backend/internal/models"

	"gorm.io/gorm"
)

type ConversationRepository interface {
	CreateConversation(tx *gorm.DB, conv *models.Conversation) error
	AddMembers(tx *gorm.DB, members []models.ConversationMember) error
	IsMember(conversationID uint, userID string) (bool, error)
	ListUserConversations(userID string) ([]models.Conversation, error)

	FindDirectConversation(userA string, userB string) (*models.Conversation, error)
}

type conversationRepository struct{ db *gorm.DB }

func NewConversationRepository(db *gorm.DB) ConversationRepository {
	return &conversationRepository{db: db}
}

func (r *conversationRepository) CreateConversation(tx *gorm.DB, conv *models.Conversation) error {
	return tx.Create(conv).Error
}

func (r *conversationRepository) AddMembers(tx *gorm.DB, members []models.ConversationMember) error {
	return tx.Create(&members).Error
}

func (r *conversationRepository) IsMember(conversationID uint, userID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.ConversationMember{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *conversationRepository) ListUserConversations(userID string) ([]models.Conversation, error) {
	var convs []models.Conversation
	err := r.db.
		Joins("JOIN conversation_members cm ON cm.conversation_id = conversations.id").
		Where("cm.user_id = ?", userID).
		Order("conversations.updated_at DESC").
		Find(&convs).Error
	return convs, err
}

func (r *conversationRepository) FindDirectConversation(userA string, userB string) (*models.Conversation, error) {
	var conv models.Conversation
	err := r.db.
		Table("conversations").
		Select("conversations.*").
		Joins("JOIN conversation_members m1 ON m1.conversation_id = conversations.id AND m1.user_id = ?", userA).
		Joins("JOIN conversation_members m2 ON m2.conversation_id = conversations.id AND m2.user_id = ?", userB).
		Where("conversations.is_group = false").
		First(&conv).Error
	if err != nil {
		return nil, err
	}
	return &conv, nil
}
