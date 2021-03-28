package v2

import "github.com/gofiber/fiber/v2"

// EndpointGetInfo handles the 'GET /v2/info' endpoint
func EndpointGetInfo(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"version": ctx.Locals("__settings_version").(string),
		"production": ctx.Locals("__settings_production").(bool),
		"settings": fiber.Map{
			"invites": ctx.Locals("__settings_invites_enabled").(bool),
		},
	})
}
