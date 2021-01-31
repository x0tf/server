package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0tf/server/internal/shared"
	"github.com/x0tf/server/internal/token"
)

// EndpointCreateInvite handles the POST /v1/invites/:code endpoint
func EndpointCreateInvite(ctx *fiber.Ctx) error {
	invites := ctx.Locals("__invites").(shared.InviteService)
	code := ctx.Params("code", token.Generate()[:32])
	if err := invites.Create(shared.Invite(code)); err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{
		"code": code,
	})
}
