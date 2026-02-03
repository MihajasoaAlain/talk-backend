package container

import (
	"time"

	"talk-backend/internal/config"
	"talk-backend/internal/http/controllers"
	"talk-backend/internal/repository"
	"talk-backend/internal/service"

	"gorm.io/gorm"
)

type App struct {
	AuthController *controllers.AuthController
	ChatController *controllers.ChatController
	UserController *controllers.UserController
}

func New(cfg *config.Config, db *gorm.DB) *App {
	userRepo := repository.NewUserRepository(db)
	rtRepo := repository.NewRefreshTokenRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	convRepo := repository.NewConversationRepository(db)
	msgRepo := repository.NewMessageRepository(db)

	authService := service.NewAuthService(
		userRepo,
		rtRepo,
		auditRepo,
		service.AuthConfig{
			JWTSecret:      cfg.JWT.Secret,
			AccessTTL:      15 * time.Minute,
			RefreshTTL:     30 * 24 * time.Hour,
			MaxFailedLogin: 5,
			LockDuration:   15 * time.Minute,
			Issuer:         "talk-backend",
		},
	)

	chatService := service.NewChatService(db, convRepo, msgRepo)
	userService := service.NewUserService(userRepo)

	// controllers
	authCtl := controllers.NewAuthController(authService)
	chatCtl := controllers.NewChatController(chatService)
	userCtl := controllers.NewUserController(userService)

	return &App{
		AuthController: authCtl,
		ChatController: chatCtl,
		UserController: userCtl,
	}
}
