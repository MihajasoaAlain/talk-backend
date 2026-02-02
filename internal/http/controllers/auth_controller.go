package controllers

import (
	"net/http"
	"talk-backend/internal/http/dto"
	"talk-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	auth *service.AuthService
}

func NewAuthController(auth *service.AuthService) *AuthController {
	return &AuthController{auth: auth}
}

func clientIP(c *gin.Context) string { return c.ClientIP() }

func userAgent(c *gin.Context) string { return c.GetHeader("User-Agent") }

// POST /auth/register
func (ctl *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctl.auth.Register(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "registered",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// POST /auth/login
func (ctl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access, refresh, err := ctl.auth.Login(req.Email, req.Password, clientIP(c), userAgent(c))
	if err != nil {
		// neutre
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

// POST /auth/refresh
func (ctl *AuthController) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access, refresh, err := ctl.auth.Refresh(req.RefreshToken, clientIP(c), userAgent(c))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (ctl *AuthController) Logout(c *gin.Context) {
	type logoutReq struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	var req logoutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_ = ctl.auth.Logout(req.RefreshToken, clientIP(c), userAgent(c))
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
