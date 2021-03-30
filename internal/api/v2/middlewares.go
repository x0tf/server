package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0tf/server/internal/shared"
	"github.com/x0tf/server/internal/token"
	"strings"
)

// MiddlewareAdminAuth handles admin token authorization
func MiddlewareAdminAuth(required bool) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Read and parse the authorization header
		header := strings.SplitN(ctx.Get(fiber.HeaderAuthorization), " ", 2)
		if len(header) != 2 || header[0] != "Bearer" {
			if required {
				return newError(fiber.StatusUnauthorized, errorCodeGenericUnauthorized, "unauthorized", nil)
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
			return newError(fiber.StatusUnauthorized, errorCodeGenericUnauthorized, "unauthorized", nil)
		}
		ctx.Locals("_is_admin", isAdmin)
		return ctx.Next()
	}
}

// MiddlewareInjectNamespace handles namespace injection and authorization
func MiddlewareInjectNamespace(handleAuth bool) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Retrieve and inject the requested namespace
		namespaces := ctx.Locals("__services_namespaces").(shared.NamespaceService)
		namespace, err := namespaces.Namespace(strings.ToLower(ctx.Params("namespace_id")))
		if err != nil {
			return err
		}
		if namespace == nil {
			if handleAuth && !ctx.Locals("_is_admin").(bool) {
				return newError(fiber.StatusUnauthorized, errorCodeGenericUnauthorized, "unauthorized", nil)
			}
			return newError(fiber.StatusNotFound, errorCodeGenericNamespaceNotFound, "namespace not found", nil)
		}
		ctx.Locals("_namespace", namespace)

		// Handle authorization if required
		if handleAuth && !ctx.Locals("_is_admin").(bool) {
			// Read and parse the authorization header
			header := strings.SplitN(ctx.Get(fiber.HeaderAuthorization), " ", 2)
			if len(header) != 2 || header[0] != "Bearer" {
				return newError(fiber.StatusUnauthorized, errorCodeGenericUnauthorized, "unauthorized", nil)
			}

			// Compare the given authentication token with the one of the found namespace
			namespace := ctx.Locals("_namespace").(*shared.Namespace)
			if valid, _ := token.Check(namespace.Token, header[1]); !valid {
				return newError(fiber.StatusUnauthorized, errorCodeGenericUnauthorized, "unauthorized", nil)
			}
		}
		return ctx.Next()
	}
}
