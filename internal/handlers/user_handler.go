package handlers

import (
	"net/http"
	"talk-backend/internal/http/response"
	"talk-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		response.InvalidBody(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": response.MsgOK,
		"user":    newUser,
	})
}
