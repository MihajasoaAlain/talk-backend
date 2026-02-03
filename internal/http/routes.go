package http

import (
	"time"

	"talk-backend/internal/container"
	"talk-backend/internal/http/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RegisterRoutes(r *gin.Engine, app *container.App, jwtSecret string) {
	loginLimiter := middleware.NewIPLimiter(rate.Every(12*time.Second), 10)

	auth := r.Group("/auth")
	{
		auth.POST("/register", app.AuthController.Register)
		auth.POST("/login", loginLimiter.Middleware(), app.AuthController.Login)
		auth.POST("/refresh", app.AuthController.Refresh)
		auth.POST("/logout", app.AuthController.Logout)
	}

	api := r.Group("/api")
	api.Use(middleware.RequireAuth(jwtSecret))
	{
		// User routes
		api.GET("/me", app.UserController.Me)

		// Conversation routes
		api.POST("/conversations/direct", app.ChatController.CreateDirect)
		api.GET("/conversations", app.ChatController.ListMyConversations)

		// Message routes
		api.POST("/conversations/:id/messages", app.ChatController.SendMessage)
		api.GET("/conversations/:id/messages", app.ChatController.GetMessages)
	}
}
