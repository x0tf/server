package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0tf/server/internal/shared"
)

// EndpointGetInvites handles the 'GET /v2/invites?limit=0&skip=10' endpoint
func EndpointGetInvites(ctx *fiber.Ctx) error {
	// Extract required services
	invites := ctx.Locals("__services_invites").(shared.InviteService)

	// Retrieve the desired amount of invites
	limit, err := parseQueryInt(ctx, "limit", 10)
	if err != nil {
		return err
	}
	skip, err := parseQueryInt(ctx, "skip", 0)
	if err != nil {
		return err
	}
	found, err := invites.Invites(limit, skip)
	if err != nil {
		return err
	}

	// Retrieve the total amount of invites
	count, err := invites.Count()
	if err != nil {
		return err
	}

	// Respond with the found invites
	return ctx.JSON(paginatedResponse{
		Data: found,
		Pagination: pagination{
			TotalElements:     count,
			DisplayedElements: len(found),
		},
	})
}