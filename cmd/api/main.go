package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vern/skillflow/internal/api/router"
	"github.com/vern/skillflow/internal/config"
	"github.com/vern/skillflow/internal/repository"
	"github.com/vern/skillflow/internal/service"
	"github.com/vern/skillflow/pkg/cache"
	"github.com/vern/skillflow/pkg/database"
	"github.com/vern/skillflow/pkg/logger"
)

func main() {
	// Initialize logger
	log := logger.New()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration", "error", err)
	}

	// Initialize database
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer database.Close(db)

	// Initialize Redis cache
	redisClient, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatal("Failed to connect to Redis", "error", err)
	}
	defer redisClient.Close()

	// Initialize repositories
	repos := repository.NewRepositories(db)

	// Initialize services
	services := service.NewServices(service.ServicesDeps{
		Repos:  repos,
		Cache:  redisClient,
		Config: cfg,
		Logger: log,
	})

	// Set Gin mode
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	r := router.NewRouter(services, cfg, log)

	// Create HTTP server
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        r,
		ReadTimeout:    time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in goroutine
	go func() {
		log.Info("Starting server", "port", cfg.Server.Port, "mode", cfg.Server.Mode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", "error", err)
	}

	log.Info("Server exited")
}
