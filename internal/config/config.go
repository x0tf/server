package config

import (
	"github.com/joho/godotenv"
	"os"
	"strings"
)

// Config represents the application configuration
type Config struct {
	DatabaseDSN         string
	APIAddress          string
	GatewayAddress      string
	GatewayRootRedirect string
	Invites             bool
	AdminTokens         []string
}

// Load loads and creates a new application configuration
func Load() (*Config, bool) {
	err := godotenv.Load()
	return &Config{
		DatabaseDSN:         os.Getenv("X0_DATABASE_DSN"),
		APIAddress:          os.Getenv("X0_API_ADDRESS"),
		GatewayAddress:      os.Getenv("X0_GATEWAY_ADDRESS"),
		GatewayRootRedirect: os.Getenv("X0_GATEWAY_ROOT_REDIRECT"),
		Invites:             os.Getenv("X0_INVITES") != "",
		AdminTokens:         strings.Split(os.Getenv("X0_ADMIN_TOKENS"), ";;"),
	}, err == nil
}
