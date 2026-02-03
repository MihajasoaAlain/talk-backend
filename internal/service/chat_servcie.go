package service

import (
	"errors"
	"time"

	"talk-backend/internal/models"
	"talk-backend/internal/repository"

	"gorm.io/gorm"
)

var ErrForbidden = errors.New("forbidden")
var ErrNotFound = errors.New("not found")

type ChatService struct {
	db       *gorm.DB
	convs    repository.ConversationRepository
	messages repository.MessageRepository
}

func NewChatService(db *gorm.DB, convs repository.ConversationRepository, messages repository.MessageRepository) *ChatService {
	return &ChatService{db: db, convs: convs, messages: messages}
}

func (s *ChatService) CreateDirectConversation(me uint, other uint) (*models.Conversation, error) {
	if conv, err := s.convs.FindDirectConversation(me, other); err == nil && conv != nil {
		return conv, nil
	}

	conv := &models.Conversation{IsGroup: false}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.convs.CreateConversation(tx, conv); err != nil {
			return err
		}
		members := []models.ConversationMember{
			{ConversationID: conv.ID, UserID: me, Role: "member"},
			{ConversationID: conv.ID, UserID: other, Role: "member"},
		}
		return s.convs.AddMembers(tx, members)
	})

	if err != nil {
		return nil, err
	}
	return conv, nil
}

func (s *ChatService) ListMyConversations(me uint) ([]models.Conversation, error) {
	return s.convs.ListUserConversations(me)
}

func (s *ChatService) SendMessage(me uint, conversationID uint, content string) (*models.Message, error) {
	ok, err := s.convs.IsMember(conversationID, me)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrForbidden
	}

	msg := &models.Message{
		ConversationID: conversationID,
		SenderID:       me,
		Content:        content,
		SentAt:         time.Now(),
	}
	if err := s.messages.Create(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *ChatService) GetMessages(me uint, conversationID uint, limit int, beforeID *uint) ([]models.Message, error) {
	ok, err := s.convs.IsMember(conversationID, me)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrForbidden
	}
	return s.messages.List(conversationID, limit, beforeID)
}
