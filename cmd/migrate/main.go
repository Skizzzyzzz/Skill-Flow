package main

import (
	"fmt"
	"os"

	"github.com/vern/skillflow/internal/config"
	"github.com/vern/skillflow/internal/domain/models"
	"github.com/vern/skillflow/pkg/database"
	"github.com/vern/skillflow/pkg/logger"
)

func main() {
	log := logger.New()

	if len(os.Args) < 2 {
		log.Fatal("Usage: migrate [up|down]")
	}

	command := os.Args[1]

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration", "error", err)
	}

	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer database.Close(db)

	switch command {
	case "up":
		log.Info("Running migrations...")
		if err := runMigrations(db); err != nil {
			log.Fatal("Migration failed", "error", err)
		}
		log.Info("Migrations completed successfully")
	case "down":
		log.Info("Rolling back migrations...")
		if err := rollbackMigrations(db); err != nil {
			log.Fatal("Rollback failed", "error", err)
		}
		log.Info("Rollback completed successfully")
	default:
		log.Fatal("Unknown command", "command", command)
	}
}

func runMigrations(db *database.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Post{},
		&models.Comment{},
		&models.Reaction{},
		&models.Connection{},
		&models.Notification{},
		&models.Message{},
		&models.Group{},
		&models.GroupMember{},
		&models.Skill{},
		&models.UserSkill{},
		&models.Endorsement{},
		&models.File{},
	)
}

func rollbackMigrations(db *database.DB) error {
	models := []interface{}{
		&models.File{},
		&models.Endorsement{},
		&models.UserSkill{},
		&models.Skill{},
		&models.GroupMember{},
		&models.Group{},
		&models.Message{},
		&models.Notification{},
		&models.Connection{},
		&models.Reaction{},
		&models.Comment{},
		&models.Post{},
		&models.Profile{},
		&models.User{},
	}

	for _, model := range models {
		if err := db.Migrator().DropTable(model); err != nil {
			return fmt.Errorf("failed to drop table: %w", err)
		}
	}

	return nil
}
