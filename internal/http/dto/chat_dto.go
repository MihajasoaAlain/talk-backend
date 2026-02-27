package dto

type DirectConversationRequest struct {
	UserID string `json:"userId" binding:"required,uuid"`
}

type SendMessageRequest struct {
	Content string `json:"content" binding:"required,min=1,max=4000"`
}
