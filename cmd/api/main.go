// @title Talk Backend API
// @version 1.0
// @description API for Talk backend services.
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"
	"talk-backend/internal/config"
	"talk-backend/internal/container"
	"talk-backend/internal/db"
	"talk-backend/internal/http"

	"github.com/gin-gonic/gin"
	_ "talk-backend/docs"
)

func main() {
	cfg := config.Load()

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	gdb, err := db.ConnectPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Migration.Valided {
		log.Println("running database migrations...")
		if err := db.Migrate(gdb); err != nil {
			log.Fatalf("database migration failed: %v", err)
		}
		log.Println("database migrations completed")
	}

	app := container.New(cfg, gdb)

	r := gin.Default()

	http.RegisterRoutes(r, app, cfg.JWT.Secret)

	log.Printf("Starting server on :%s (%s)", cfg.App.Port, cfg.App.Env)
	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
