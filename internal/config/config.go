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
	Invites        bool
}

// Load loads and creates a new application configuration
func Load() (*Config, bool) {
	err := godotenv.Load()
	return &Config{
		DatabaseDSN:    os.Getenv("X0_DATABASE_DSN"),
		APIAddress:     os.Getenv("X0_API_ADDRESS"),
		GatewayAddress: os.Getenv("X0_GATEWAY_ADDRESS"),
		Invites:        os.Getenv("X0_INVITES") != "",
	}, err == nil
}
