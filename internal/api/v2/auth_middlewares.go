package v2

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

// MiddlewareAdminAuth handles admin token authorization
func MiddlewareAdminAuth(required bool) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Read and parse the authorization header
		header := strings.SplitN(ctx.Get(fiber.HeaderAuthorization), " ", 2)
		if len(header) != 2 || header[0] != "Bearer" {
			if required {
				return errUnauthorized
			}
			ctx.Locals("_is_admin", false)
			return ctx.Next()
		}

		// Check if the given token is an admin token
		isAdmin := false
		for _, adminToken := range ctx.Locals("__settings_admin_tokens").([]string) {
			if header[1] == adminToken {
				isAdmin = true
				break
			}
		}
		if required && !isAdmin {
			return errForbidden
		}
		ctx.Locals("_is_admin", isAdmin)
		return ctx.Next()
	}
}
