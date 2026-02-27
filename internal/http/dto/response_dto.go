package dto

import (
	"time"

	"talk-backend/internal/models"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type UserPublic struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegisterResponse struct {
	Message string     `json:"message"`
	User    UserPublic `json:"user"`
}

type UserMe struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type MeResponse struct {
	User UserMe `json:"user"`
}

type ConversationResponse struct {
	Conversation models.Conversation `json:"conversation"`
}

type ConversationsResponse struct {
	Conversations []models.Conversation `json:"conversations"`
}

type MessageResponseData struct {
	Message models.Message `json:"message"`
}

type MessagesResponse struct {
	Messages []models.Message `json:"messages"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
