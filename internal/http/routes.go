package http

import (
	"talk-backend/internal/container"
	"talk-backend/internal/http/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RegisterRoutes(r *gin.Engine, app *container.App) {
	loginLimiter := middleware.NewIPLimiter(rate.Every(12*time.Second), 10)

	auth := r.Group("/auth")
	{
		auth.POST("/register", app.AuthController.Register)
		auth.POST("/login", loginLimiter.Middleware(), app.AuthController.Login)
		auth.POST("/refresh", app.AuthController.Refresh)
		auth.POST("/logout", app.AuthController.Logout)
	}
}
