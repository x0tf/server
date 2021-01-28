package config

import (
	"github.com/joho/godotenv"
	"os"
)

// Config represents the application configuration
type Config struct {
	DatabaseDSN    string
	APIAddress     string
	GatewayAddress string
}

// Load loads and creates a new application configuration
func Load() (*Config, bool) {
	err := godotenv.Load()
	return &Config{
		DatabaseDSN:    os.Getenv("X0_DATABASE_DSN"),
		APIAddress:     os.Getenv("X0_API_ADDRESS"),
		GatewayAddress: os.Getenv("X0_GATEWAY_ADDRESS"),
	}, err == nil
}
