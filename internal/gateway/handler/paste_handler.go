package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0tf/server/internal/shared"
)

// pasteHandler handles incoming requests for paste elements
func pasteHandler(ctx *fiber.Ctx) error {
	element := ctx.Locals("_element").(*shared.Element)
	return ctx.SendString(element.PublicData["content"].(string))
}
