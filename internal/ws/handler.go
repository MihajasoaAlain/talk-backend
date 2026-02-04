package ws

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"talk-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

type WSHandler struct {
	hub       *Hub
	chat      *service.ChatService
	jwtSecret string
}

func NewWSHandler(hub *Hub, chat *service.ChatService, jwtSecret string) *WSHandler {
	return &WSHandler{hub: hub, chat: chat, jwtSecret: jwtSecret}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type InMessage struct {
	Type           string `json:"type"`
	ConversationID uint   `json:"conversationId"`
	Content        string `json:"content"`
}

func (h *WSHandler) Handle(c *gin.Context) {
	convStr := c.Query("conversationId")
	conv64, err := strconv.ParseUint(convStr, 10, 64)
	if err != nil || conv64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "conversationId required"})
		return
	}
	roomID := uint(conv64)

	userID, ok := h.extractUserID(c.GetHeader("Authorization"))
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	_, checkErr := h.chat.GetMessages(userID, roomID, 1, nil)
	if checkErr == service.ErrForbidden {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	if checkErr != nil && checkErr != service.ErrForbidden {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &Client{
		conn:   conn,
		hub:    h.hub,
		send:   make(chan []byte, 64),
		userID: userID,
		roomID: roomID,
	}
	h.hub.register <- client

	go client.writePump()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var in InMessage
		if err := json.Unmarshal(p, &in); err != nil {
			continue
		}
		if in.Type != "message" || in.Content == "" || in.ConversationID != roomID {
			continue
		}

		msg, err := h.chat.SendMessage(userID, roomID, in.Content)
		if err != nil {
			continue
		}

		out, _ := json.Marshal(gin.H{
			"type": "message",
			"message": gin.H{
				"id":             msg.ID,
				"conversationId": msg.ConversationID,
				"senderId":       msg.SenderID,
				"content":        msg.Content,
				"sentAt":         msg.SentAt,
			},
		})
		h.hub.broadcast <- RoomMessage{RoomID: roomID, Data: out}
	}

	h.hub.unregister <- client
	_ = conn.Close()
}

func (h *WSHandler) extractUserID(authHeader string) (uint, bool) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
		return 0, false
	}

	tok, err := jwt.Parse(parts[1], func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(h.jwtSecret), nil
	})
	if err != nil || tok == nil || !tok.Valid {
		return 0, false
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false
	}

	sub, ok := claims["sub"]
	if !ok {
		return 0, false
	}

	switch v := sub.(type) {
	case float64:
		return uint(v), true
	case uint:
		return v, true
	case string:
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil || n == 0 {
			return 0, false
		}
		return uint(n), true
	default:
		return 0, false
	}
}
