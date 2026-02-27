package ws

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"talk-backend/internal/http/response"
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
	CheckOrigin: func(r *http.Request) bool { return true }, // en prod: restreindre
}

var uuidV4LikeRe = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[1-5][a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`)

type InMessage struct {
	Type           string `json:"type"`
	ConversationID uint   `json:"conversationId"`
	Content        string `json:"content,omitempty"`
	IsTyping       *bool  `json:"isTyping,omitempty"`
}

func (h *WSHandler) Handle(c *gin.Context) {
	convStr := c.Query("conversationId")
	conv64, err := strconv.ParseUint(convStr, 10, 64)
	if err != nil || conv64 == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidRequest, response.MsgConversationRequired)
		return
	}
	roomID := uint(conv64)

	userID, ok := h.extractUserID(c.GetHeader("Authorization"))
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, response.MsgUnauthorized)
		return
	}

	_, checkErr := h.chat.GetMessages(userID, roomID, 1, nil)
	if checkErr == service.ErrForbidden {
		response.Error(c, http.StatusForbidden, response.CodeForbidden, response.MsgForbidden)
		return
	}
	if checkErr != nil && checkErr != service.ErrForbidden {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, response.MsgInternalServer)
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

	lastTyping := time.Time{}

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var in InMessage
		if err := json.Unmarshal(p, &in); err != nil {
			continue
		}

		if in.ConversationID != roomID {
			continue
		}

		if in.Type == "typing" && in.IsTyping != nil {
			if time.Since(lastTyping) < 500*time.Millisecond {
				continue
			}
			lastTyping = time.Now()

			outTyping, _ := json.Marshal(gin.H{
				"type":           "typing",
				"conversationId": roomID,
				"userId":         userID,
				"isTyping":       *in.IsTyping,
			})

			h.hub.broadcast <- RoomMessage{RoomID: roomID, Data: outTyping}
			continue
		}

		if in.Type != "message" || strings.TrimSpace(in.Content) == "" {
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

func (h *WSHandler) extractUserID(authHeader string) (string, bool) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
		return "", false
	}

	tok, err := jwt.Parse(parts[1], func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(h.jwtSecret), nil
	})
	if err != nil || tok == nil || !tok.Valid {
		return "", false
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}

	sub, ok := claims["sub"]
	if !ok {
		return "", false
	}

	switch v := sub.(type) {
	case string:
		if !isUUID(v) {
			return "", false
		}
		return v, true
	default:
		return "", false
	}
}

func isUUID(v string) bool {
	return uuidV4LikeRe.MatchString(v)
}
