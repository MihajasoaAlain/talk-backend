package main

import (
	"log"
	"talk-backend/internal/config"
	"talk-backend/internal/db"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	gdb, err := db.ConnectPostgres(cfg)
	log.Println(gdb)
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	if cfg.Migration.Valided {
		log.Println("running database migrations...")
		if err := db.Migrate(gdb); err != nil {
			log.Fatalf("database migration failed: %v", err)
		}
		log.Println("database migrations completed")
	}

	log.Printf("Starting server on :%s (%s)", cfg.App.Port, cfg.App.Env)
	r.Run(":" + cfg.App.Port)
}
