package main

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
	"github.com/mdeadwiler/EyeSky/internal/platform/config"
	"github.com/mdeadwiler/EyeSky/internal/platform/db"
	"github.com/mdeadwiler/EyeSky/internal/platform/logging"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Failed to load .env file: %v\n", err)
	}

	// Load configuration
	cfg, err := config.New()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Debug: Print database config
	fmt.Printf("DB Config - User: %s, DB: %s, Host: %s\n", cfg.Database.User, cfg.Database.DBName, cfg.Database.Host)


	
	// Initialize logger
	logger := logging.New(cfg.Logging)
	logger.Info("Starting EyeSky Gateway Service")

	// Initialize database
	database, err := db.New(cfg.Database)
	if err != nil {
		logger.ErrorWithFields("Failed to connect to database", map[string]interface{}{
			"error": err.Error(),
		})
		os.Exit(1)
	}
	defer database.Close()

	// Run database migrations
	if err := database.Migrate(); err != nil {
		logger.ErrorWithFields("Failed to run migrations", map[string]interface{}{
			"error": err.Error(),
		})
		os.Exit(1)
	}

	logger.InfoWithFields("EyeSky Gateway Service started successfully", map[string]interface{}{
		"environment": cfg.Environment,
		"port":        cfg.Server.Port,
	})

	// TODO: Initialize repositories, services, and HTTP server
	select {} // Keep the service running
}
