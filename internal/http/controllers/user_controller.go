package controllers

import (
	"net/http"

	"talk-backend/internal/http/dto"
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

// Me godoc
// @Summary Get current user
// @Description Return the authenticated user's profile.
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.MeResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/me [get]
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

	c.JSON(http.StatusOK, dto.MeResponse{
		User: dto.UserMe{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}
