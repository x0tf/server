package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0tf/server/internal/shared"
	"github.com/x0tf/server/internal/utils"
	"time"
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

// EndpointGetInvite handles the 'GET /v2/invites/:invite_code' endpoint
func EndpointGetInvite(ctx *fiber.Ctx) error {
	return ctx.JSON(ctx.Locals("_invite").(*shared.Invite))
}

type endpointCreateInviteRequestBody struct {
	Code    *string `json:"code" xml:"code" form:"code"`
	MaxUses *int    `json:"max_uses" xml:"max_uses" form:"max_uses"`
}

// EndpointCreateInvite handles the 'POST /v2/invites' endpoint
func EndpointCreateInvite(ctx *fiber.Ctx) error {
	// Extract required services
	invites := ctx.Locals("__services_invites").(shared.InviteService)

	// Try to parse the body into a request body struct
	body := new(endpointCreateInviteRequestBody)
	if err := ctx.BodyParser(body); err != nil {
		return errorGenericBadRequestBody
	}

	// Validate the wished code if one was provided or generate a new one
	var code string
	if body.Code != nil {
		wishedCode := *body.Code
		found, err := invites.Invite(wishedCode)
		if err != nil {
			return err
		}
		if found != nil {
			return errorInviteInviteCodeInUse
		}
		code = wishedCode
	} else {
		for {
			generated := utils.GenerateInviteCode()
			found, err := invites.Invite(generated)
			if err != nil {
				return err
			}
			if found == nil {
				code = generated
				break
			}
		}
	}

	// Define the maximum amount of uses
	maxUses := -1
	if body.MaxUses != nil {
		maxUses = *body.MaxUses
	}

	// Create and respond with the invite
	invite := &shared.Invite{
		Code:    code,
		Uses:    0,
		MaxUses: maxUses,
		Created: time.Now().Unix(),
	}
	if err := invites.CreateOrReplace(invite); err != nil {
		return err
	}
	return ctx.JSON(invite)
}

type endpointPatchInviteRequestBody struct {
	Code    *string `json:"code" xml:"code" form:"code"`
	MaxUses *int    `json:"max_uses" xml:"max_uses" form:"max_uses"`
}

// EndpointPatchInvite handles the 'PATCH /v2/invites/:invite_code' endpoint
func EndpointPatchInvite(ctx *fiber.Ctx) error {
	// Extract required services
	invites := ctx.Locals("__services_invites").(shared.InviteService)

	// Extract required resources
	invite := ctx.Locals("_invite").(*shared.Invite)

	// Try to parse the body into a request body struct
	body := new(endpointPatchInviteRequestBody)
	if err := ctx.BodyParser(body); err != nil {
		return errorGenericBadRequestBody
	}

	// Validate and update the code of the invite if specified
	code := invite.Code
	if body.Code != nil {
		code = *body.Code
		found, err := invites.Invite(code)
		if err != nil {
			return err
		}
		if found != nil {
			return errorInviteInviteCodeInUse
		}
	}

	// Update the maximum amount of uses if specified
	if body.MaxUses != nil {
		invite.MaxUses = *body.MaxUses
	}

	// Update the invite and respond with the updated version of it
	if code != invite.Code {
		if err := invites.Delete(invite.Code); err != nil {
			return err
		}
		invite.Code = code
	}
	if err := invites.CreateOrReplace(invite); err != nil {
		return err
	}
	return ctx.JSON(invite)
}
