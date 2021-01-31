package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0tf/server/internal/shared"
	"github.com/x0tf/server/internal/token"
	"strings"
)

// MiddlewareTokenAuth handles namespace token authentication
func MiddlewareTokenAuth(ctx *fiber.Ctx) error {
	// Extract the namespace service, retrieve the requested namespace and check if it exists
	namespaces := ctx.Locals("__namespaces").(shared.NamespaceService)
	namespace, err := namespaces.Namespace(ctx.Params("namespace"))
	if err != nil {
		return err
	}
	if namespace == nil {
		return fiber.NewError(fiber.StatusNotFound, "that namespace does not exist")
	}

	// Perform user authentication if the request was not made by an admin
	isAdmin, _ := ctx.Locals("_admin").(bool)
	if !isAdmin {
		// Read and validate the header itself
		header := strings.SplitN(ctx.Get(fiber.HeaderAuthorization), " ", 2)
		if len(header) != 2 || header[0] != "Bearer" {
			return fiber.ErrUnauthorized
		}

		// Compare the given authentication token with the one of the found namespace
		if valid, _ := token.Check(namespace.Token, header[1]); !valid {
			return fiber.ErrUnauthorized
		}
	}

	// Inject the namespace this request is aimed at
	ctx.Locals("_namespace", namespace)
	return ctx.Next()
}

// MiddlewareAdminAuth handles admin token authentication
func MiddlewareAdminAuth(ctx *fiber.Ctx) error {
	header := strings.SplitN(ctx.Get(fiber.HeaderAuthorization), " ", 2)
	if len(header) != 2 || header[0] != "Bearer" {
		ctx.Locals("_admin", false)
		return ctx.Next()
	}

	isAdmin := false
	for _, adminToken := range ctx.Locals("__admin_tokens").([]string) {
		if header[1] == adminToken {
			isAdmin = true
			break
		}
	}
	ctx.Locals("_admin", isAdmin)
	return ctx.Next()
}

// MiddlewareRequireAdminAuth handles admin token authentication requirement
func MiddlewareRequireAdminAuth(ctx *fiber.Ctx) error {
	if !ctx.Locals("_admin").(bool) {
		return fiber.ErrForbidden
	}
	return ctx.Next()
}
