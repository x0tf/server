package config

import (
	"github.com/joho/godotenv"
	"github.com/x0tf/server/internal/env"
)

// Config represents the application configuration
type Config struct {
	DatabaseDSN          string
	APIAddress           string
	APIRequestsPerMinute int
	GatewayAddress       string
	GatewayRootRedirect  string
	Invites              bool
	AdminTokens          []string
}

// Load loads and creates a new application configuration
func Load() (*Config, bool) {
	err := godotenv.Load()
	return &Config{
		DatabaseDSN:          env.MustString("X0_DATABASE_DSN", ""),
		APIAddress:           env.MustString("X0_API_ADDRESS", ":8080"),
		APIRequestsPerMinute: env.MustInt("X0_API_REQUESTS_PER_MINUTE", 60),
		GatewayAddress:       env.MustString("X0_GATEWAY_ADDRESS", ":8081"),
		GatewayRootRedirect:  env.MustString("X0_GATEWAY_ROOT_REDIRECT", ""),
		Invites:              env.MustBool("X0_INVITES", false),
		AdminTokens:          env.MustStringSlice("X0_ADMIN_TOKENS", ";;", []string{}),
	}, err == nil
}
