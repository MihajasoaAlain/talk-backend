package controllers

import (
	"errors"
	"net/http"
	"talk-backend/internal/http/dto"
	"talk-backend/internal/http/response"
	"talk-backend/internal/repository"
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

// Register godoc
// @Summary Register a new user
// @Description Create a new user account.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register payload"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (ctl *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidBody(c, err)
		return
	}

	user, err := ctl.auth.Register(req.Username, req.Email, req.Password, req.AvatarURL)
	if err != nil {
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			response.Error(c, http.StatusConflict, response.CodeEmailAlreadyExists, response.MsgEmailAlreadyExists)
			return
		}
		response.Error(c, http.StatusInternalServerError, response.CodeRegisterFailed, response.MsgRegisterFailed)
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterResponse{
		Message: response.MsgRegistered,
		User: dto.UserPublic{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
		},
	})
}

// POST /auth/login
// Login godoc
// @Summary Login
// @Description Authenticate a user and return access and refresh tokens.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login payload"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (ctl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidBody(c, err)
		return
	}

	access, refresh, err := ctl.auth.Login(req.Email, req.Password, clientIP(c), userAgent(c))
	if err != nil {
		response.Error(c, http.StatusUnauthorized, response.CodeInvalidCredentials, response.MsgInvalidCredentials)
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

// POST /auth/refresh
// Refresh godoc
// @Summary Refresh tokens
// @Description Refresh access and refresh tokens using a refresh token.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh payload"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/refresh [post]
func (ctl *AuthController) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidBody(c, err)
		return
	}

	access, refresh, err := ctl.auth.Refresh(req.RefreshToken, clientIP(c), userAgent(c))
	if err != nil {
		response.Error(c, http.StatusUnauthorized, response.CodeInvalidRefreshToken, response.MsgInvalidRefreshToken)
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

// Logout godoc
// @Summary Logout
// @Description Revoke a refresh token.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LogoutRequest true "Logout payload"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /auth/logout [post]
func (ctl *AuthController) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidBody(c, err)
		return
	}

	_ = ctl.auth.Logout(req.RefreshToken, clientIP(c), userAgent(c))
	c.JSON(http.StatusOK, dto.MessageResponse{Message: response.MsgOK})
}
