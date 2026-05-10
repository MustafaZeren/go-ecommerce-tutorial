package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/mustafazeren/go-ecommerce-course/internal/config"
	"github.com/mustafazeren/go-ecommerce-course/internal/database"
	"github.com/mustafazeren/go-ecommerce-course/internal/logger"
)

func main() {
	log := logger.New()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}
	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	mainDb, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get DB connection")
	}
	defer func(mainDb *sql.DB) {
		err := mainDb.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to close DB connection")
		}
	}(mainDb)
	gin.SetMode(cfg.Server.GinMode)
	log.Info().Msg("Start server")
}
