package dto

type DirectConversationRequest struct {
	UserID uint `json:"userId" binding:"required"`
}

type SendMessageRequest struct {
	Content string `json:"content" binding:"required,min=1,max=4000"`
}
