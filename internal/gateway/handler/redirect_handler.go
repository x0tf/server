package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0tf/server/internal/shared"
)

// pasteHandler handles incoming requests for redirect elements
func redirectHandler(ctx *fiber.Ctx) error {
	element := ctx.Locals("_element").(*shared.Element)
	return ctx.Redirect(element.PublicData["target_url"].(string), fiber.StatusTemporaryRedirect)
}
