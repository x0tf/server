package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0tf/server/internal/shared"
	"github.com/x0tf/server/internal/token"
)

// EndpointListInvites handles the GET /v1/invites endpoint
func EndpointListInvites(ctx *fiber.Ctx) error {
	invites := ctx.Locals("__invites").(shared.InviteService)
	list, err := invites.Invites()
	if err != nil {
		return err
	}
	return ctx.JSON(list)
}

// EndpointValidateInvite handles the GET /v1/invites/:code endpoint
func EndpointValidateInvite(ctx *fiber.Ctx) error {
	code := ctx.Params("code")
	invites := ctx.Locals("__invites").(shared.InviteService)
	valid, err := invites.IsValid(code)
	if err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{
		"valid": valid,
	})
}

// EndpointCreateInvite handles the POST /v1/invites/:code? endpoint
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

// EndpointDeleteInvite handles the DELETE /v1/invites/:code endpoint
func EndpointDeleteInvite(ctx *fiber.Ctx) error {
	code := ctx.Params("code")
	invites := ctx.Locals("__invites").(shared.InviteService)
	return invites.Delete(code)
}
