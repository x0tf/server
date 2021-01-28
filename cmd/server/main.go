package main

import (
	"github.com/x0tf/server/internal/config"
	"github.com/x0tf/server/internal/database/postgres"
	"log"
)

func main() {
	// Initialize the application configuration
	cfg, usedDotEnv := config.Load()
	if !usedDotEnv {
		log.Println("NOTE: No .env file was found. This is no error and the application will use the systems environment variables.")
	}

	// Initialize the namespace service
	namespaces, err := postgres.NewNamespaceService(cfg.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	// Initialize the element service
	elements, err := postgres.NewElementService(cfg.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	// TODO: Implement startup logic
}
