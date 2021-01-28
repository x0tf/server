package v1

import (
	"github.com/gofiber/fiber/v2"
)

// EndpointGetInfo handles the GET /v1/info endpoint
func EndpointGetInfo(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"production": ctx.Locals("__production").(bool),
		"version":    ctx.Locals("__version").(string),
	})
}
