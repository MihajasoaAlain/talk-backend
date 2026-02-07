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
	"time"

	_ "talk-backend/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	http.RegisterRoutes(r, app, cfg.JWT.Secret)

	log.Printf("Starting server on :%s (%s)", cfg.App.Port, cfg.App.Env)
	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
