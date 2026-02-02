package container

import (
	"talk-backend/internal/config"
	"talk-backend/internal/http/controllers"
	"talk-backend/internal/repository"
	"talk-backend/internal/service"

	"gorm.io/gorm"
)

type App struct {
	AuthController *controllers.AuthController
}

func New(cfg *config.Config, db *gorm.DB) *App {
	userRepo := repository.NewUserRepository(db)
	rtRepo := repository.NewRefreshTokenRepository(db)
	auditRepo := repository.NewAuditRepository(db)

	authService := service.NewAuthService(
		userRepo,
		rtRepo,
		auditRepo,
		service.AuthConfig{
			JWTSecret:      cfg.JWT.Secret,
			AccessTTL:      15 * 60, // si ton AuthConfig utilise time.Duration: remplace par 15*time.Minute
			RefreshTTL:     30 * 24 * 60 * 60,
			MaxFailedLogin: 5,
			LockDuration:   15 * 60,
			Issuer:         "talk-backend",
		},
	)

	// controllers
	authCtl := controllers.NewAuthController(authService)

	return &App{
		AuthController: authCtl,
	}
}
