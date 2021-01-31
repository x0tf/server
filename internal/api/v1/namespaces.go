package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0tf/server/internal/shared"
	"github.com/x0tf/server/internal/token"
	"github.com/x0tf/server/internal/validation"
	"strings"
)

// MiddlewareTokenAuth handles namespace token authentication
func MiddlewareTokenAuth(ctx *fiber.Ctx) error {
	// Read and validate the header itself
	header := strings.SplitN(ctx.Get(fiber.HeaderAuthorization), " ", 1)
	if len(header) != 2 || header[0] != "Bearer" {
		return fiber.ErrUnauthorized
	}

	// Extract the namespace service, retrieve the requested namespace and check if it exists
	namespaces := ctx.Locals("__namespaces").(shared.NamespaceService)
	namespace, err := namespaces.Namespace(ctx.Params("namespace"))
	if err != nil {
		return err
	}
	if namespace == nil {
		return fiber.ErrUnauthorized
	}

	// Compare the given authentication token with the one of the found namespace
	if valid, _ := token.Check(namespace.Token, header[1]); !valid {
		return fiber.ErrUnauthorized
	}

	// Inject the namespace this request is aimed at
	ctx.Locals("_namespace", namespace)
	return ctx.Next()
}

// EndpointCreateNamespace handles the POST /v1/namespaces/:namespace endpoint
func EndpointCreateNamespace(ctx *fiber.Ctx) error {
	// Validate the given namespace ID
	id := ctx.Params("namespace")
	if errors := validation.ValidateNamespaceID(id); len(errors) > 0 {
		messages := make([]string, 0, len(errors))
		for _, err := range errors {
			messages = append(messages, err.Error())
		}
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(messages)
	}

	// Check if a namespace with this ID already exists
	namespaces := ctx.Locals("__namespaces").(shared.NamespaceService)
	found, err := namespaces.Namespace(id)
	if err != nil {
		return err
	}
	if found != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "the given namespace ID is already taken")
	}

	// Create a new namespace and a copy of it
	namespace := &shared.Namespace{
		ID:     id,
		Token:  token.Generate(),
		Active: true,
	}
	namespaceCopy := *namespace

	// Hash the token of the original namespace and insert it into the database
	hash, err := token.Hash(namespace.Token)
	if err != nil {
		return err
	}
	namespace.Token = hash
	if err = namespaces.CreateOrReplace(namespace); err != nil {
		return err
	}

	// Return the copied namespace with the raw token still placed in it
	return ctx.JSON(namespaceCopy)
}
