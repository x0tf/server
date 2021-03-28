package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0tf/server/internal/shared"
)

// EndpointGetNamespaces handles the 'GET /v2/namespaces?limit=0&skip=10' endpoint
func EndpointGetNamespaces(ctx *fiber.Ctx) error {
	// Extract required services
	namespaces := ctx.Locals("__services_namespaces").(shared.NamespaceService)

	// Retrieve the desired amount of namespaces
	limit, err := parseQueryInt(ctx, "limit", 10)
	if err != nil {
		return err
	}
	skip, err := parseQueryInt(ctx, "skip", 0)
	if err != nil {
		return err
	}
	found, err := namespaces.Namespaces(limit, skip)
	if err != nil {
		return err
	}

	// Retrieve the total amount of namespaces
	count, err := namespaces.Count()
	if err != nil {
		return err
	}

	// Remove the tokens of these namespaces
	processed := make([]shared.Namespace, 0, len(found))
	for _, namespace := range found {
		tmp := *namespace
		tmp.Token = ""
		processed = append(processed, tmp)
	}

	// Respond with the processed namespaces
	return ctx.JSON(paginatedResponse{
		Data: processed,
		Pagination: pagination{
			TotalElements:     count,
			DisplayedElements: len(processed),
		},
	})
}

// EndpointGetNamespace handles the 'GET /v2/namespaces/:namespace_id' endpoint
func EndpointGetNamespace(ctx *fiber.Ctx) error {
	namespace := *(ctx.Locals("_namespace").(*shared.Namespace))
	namespace.Token = ""
	return ctx.JSON(namespace)
}
