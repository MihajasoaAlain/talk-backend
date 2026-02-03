package controllers

import (
	"net/http"
	"strconv"

	"talk-backend/internal/http/dto"
	"talk-backend/internal/http/middleware"
	"talk-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	chat *service.ChatService
}

func NewChatController(chat *service.ChatService) *ChatController {
	return &ChatController{chat: chat}
}

func (ctl *ChatController) CreateDirect(c *gin.Context) {
	me, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.DirectConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conv, err := ctl.chat.CreateDirectConversation(me, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create conversation"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"conversation": conv})
}

func (ctl *ChatController) ListMyConversations(c *gin.Context) {
	me, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	convs, err := ctl.chat.ListMyConversations(me)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list conversations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"conversations": convs})
}

func (ctl *ChatController) SendMessage(c *gin.Context) {
	me, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	convID64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || convID64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}
	convID := uint(convID64)

	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, err := ctl.chat.SendMessage(me, convID, req.Content)
	if err != nil {
		if err == service.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send message"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": msg})
}

func (ctl *ChatController) GetMessages(c *gin.Context) {
	me, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	convID64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || convID64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
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
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": msgs})
}
