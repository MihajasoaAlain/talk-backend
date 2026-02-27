package controllers

import (
	"net/http"
	"strconv"

	"talk-backend/internal/http/dto"
	"talk-backend/internal/http/middleware"
	"talk-backend/internal/http/response"
	"talk-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	chat *service.ChatService
}

func NewChatController(chat *service.ChatService) *ChatController {
	return &ChatController{chat: chat}
}

// CreateDirect godoc
// @Summary Create a direct conversation
// @Description Create a one-to-one conversation with another user.
// @Tags conversations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.DirectConversationRequest true "Direct conversation payload"
// @Success 201 {object} dto.ConversationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/conversations/direct [post]
func (ctl *ChatController) CreateDirect(c *gin.Context) {
	me, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, response.MsgUnauthorized)
		return
	}

	var req dto.DirectConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidBody(c, err)
		return
	}

	conv, err := ctl.chat.CreateDirectConversation(me, req.UserID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeConversationFailed, response.MsgCreateConversation)
		return
	}

	c.JSON(http.StatusCreated, dto.ConversationResponse{Conversation: *conv})
}

// ListMyConversations godoc
// @Summary List my conversations
// @Description Return conversations the authenticated user is a member of.
// @Tags conversations
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.ConversationsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/conversations [get]
func (ctl *ChatController) ListMyConversations(c *gin.Context) {
	me, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, response.MsgUnauthorized)
		return
	}

	convs, err := ctl.chat.ListMyConversations(me)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeConversationFailed, response.MsgListConversations)
		return
	}

	c.JSON(http.StatusOK, dto.ConversationsResponse{Conversations: convs})
}

// SendMessage godoc
// @Summary Send a message
// @Description Send a message in a conversation.
// @Tags messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Conversation ID"
// @Param request body dto.SendMessageRequest true "Message payload"
// @Success 201 {object} dto.MessageResponseData
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/conversations/{id}/messages [post]
func (ctl *ChatController) SendMessage(c *gin.Context) {
	me, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, response.MsgUnauthorized)
		return
	}

	convID64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || convID64 == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidRequest, response.MsgInvalidConversation)
		return
	}
	convID := uint(convID64)

	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidBody(c, err)
		return
	}

	msg, err := ctl.chat.SendMessage(me, convID, req.Content)
	if err != nil {
		if err == service.ErrForbidden {
			response.Error(c, http.StatusForbidden, response.CodeForbidden, response.MsgForbidden)
			return
		}
		response.Error(c, http.StatusInternalServerError, response.CodeMessageFailed, response.MsgSendMessage)
		return
	}

	c.JSON(http.StatusCreated, dto.MessageResponseData{Message: *msg})
}

// GetMessages godoc
// @Summary Get messages
// @Description Get messages in a conversation.
// @Tags messages
// @Security BearerAuth
// @Produce json
// @Param id path int true "Conversation ID"
// @Param limit query int false "Max messages to return"
// @Param beforeId query int false "Return messages before this message ID"
// @Success 200 {object} dto.MessagesResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/conversations/{id}/messages [get]
func (ctl *ChatController) GetMessages(c *gin.Context) {
	me, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, response.MsgUnauthorized)
		return
	}

	convID64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || convID64 == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidRequest, response.MsgInvalidConversation)
		return
	}
	convID := uint(convID64)

	limit, _ := strconv.Atoi(c.Query("limit"))
	var beforeID *uint
	if v := c.Query("beforeId"); v != "" {
		b, err := strconv.ParseUint(v, 10, 64)
		if err == nil && b > 0 {
			tmp := uint(b)
			beforeID = &tmp
		}
	}

	msgs, err := ctl.chat.GetMessages(me, convID, limit, beforeID)
	if err != nil {
		if err == service.ErrForbidden {
			response.Error(c, http.StatusForbidden, response.CodeForbidden, response.MsgForbidden)
			return
		}
		response.Error(c, http.StatusInternalServerError, response.CodeMessageFailed, response.MsgGetMessages)
		return
	}

	c.JSON(http.StatusOK, dto.MessagesResponse{Messages: msgs})
}
