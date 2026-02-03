package controllers

import (
	"net/http"

	"talk-backend/internal/http/middleware"
	"talk-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	user *service.UserService
}

func NewUserController(user *service.UserService) *UserController {
	return &UserController{user: user}
}

func (ctl *UserController) Me(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := ctl.user.GetMe(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		},
	})
}
