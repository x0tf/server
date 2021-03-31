package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0tf/server/internal/validation"
)

// EndpointGetInfo handles the 'GET /v2/info' endpoint
func EndpointGetInfo(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"version":    ctx.Locals("__settings_version").(string),
		"production": ctx.Locals("__settings_production").(bool),
		"settings": fiber.Map{
			"invites": ctx.Locals("__settings_invites_enabled").(bool),
			"namespace_id_rules": fiber.Map{
				"min_length":         validation.NamespaceIDMinimumLength,
				"max_length":         validation.NamespaceIDMaximumLength,
				"allowed_characters": validation.NamespaceIDAllowedCharacters,
			},
		},
	})
}
