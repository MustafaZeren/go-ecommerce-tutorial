// Starting point of this project.
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mustafazeren/go-ecommerce-course/internal/config"
	"github.com/mustafazeren/go-ecommerce-course/internal/database"
	"github.com/mustafazeren/go-ecommerce-course/internal/logger"
	"github.com/mustafazeren/go-ecommerce-course/internal/server"
)

func main() {

	log := logger.New()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error().Err(err).Msg("Failed to load config")
		return
	}

	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		return
	}

	mainDb, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get DB connection")
		return
	}

	defer func(mainDb *sql.DB) {
		if err := mainDb.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close DB connection")
		}
	}(mainDb)

	gin.SetMode(cfg.Server.GinMode)

	srv := server.New(cfg, db, log)
	router := srv.SetupRoutes()
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("Starting Http Server...")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Failed to start Http Server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to shutdown Http Server")
		return
	}

	log.Info().Msg("Server gracefully stopped.")
}
