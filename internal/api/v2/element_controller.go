package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0tf/server/internal/shared"
)

// EndpointGetElements handles the 'GET /v2/elements?limit=0&skip=10' endpoint
func EndpointGetElements(ctx *fiber.Ctx) error {
	// Extract required services
	elements := ctx.Locals("__services_elements").(shared.ElementService)

	// Retrieve the desired amount of elements
	limit, err := parseQueryInt(ctx, "limit", 10)
	if err != nil {
		return err
	}
	skip, err := parseQueryInt(ctx, "skip", 0)
	if err != nil {
		return err
	}
	found, err := elements.Elements(limit, skip)
	if err != nil {
		return err
	}

	// Retrieve the total amount of elements
	count, err := elements.Count()
	if err != nil {
		return err
	}

	// Remove the internal data fields of these elements
	processed := make([]shared.Element, 0, len(found))
	for _, element := range found {
		tmp := *element
		tmp.InternalData = map[string]interface{}{}
		processed = append(processed, tmp)
	}

	// Respond with the found elements
	return ctx.JSON(paginatedResponse{
		Data: processed,
		Pagination: pagination{
			TotalElements:     count,
			DisplayedElements: len(processed),
		},
	})
}

// EndpointGetNamespaceElements handles the 'GET /v2/elements/:namespace_id?limit=0&skip=10' endpoint
func EndpointGetNamespaceElements(ctx *fiber.Ctx) error {
	// Extract required services
	elements := ctx.Locals("__services_elements").(shared.ElementService)

	// Extract required resources
	namespace := ctx.Locals("_namespace").(*shared.Namespace)

	// Retrieve the desired amount of elements
	limit, err := parseQueryInt(ctx, "limit", 10)
	if err != nil {
		return err
	}
	skip, err := parseQueryInt(ctx, "skip", 0)
	if err != nil {
		return err
	}
	found, err := elements.ElementsInNamespace(namespace.ID, limit, skip)
	if err != nil {
		return err
	}

	// Retrieve the total amount of elements
	count, err := elements.CountInNamespace(namespace.ID)
	if err != nil {
		return err
	}

	// Remove the internal data fields of these elements
	processed := make([]shared.Element, 0, len(found))
	for _, element := range found {
		tmp := *element
		tmp.InternalData = map[string]interface{}{}
		processed = append(processed, tmp)
	}

	// Respond with the found elements
	return ctx.JSON(paginatedResponse{
		Data: processed,
		Pagination: pagination{
			TotalElements:     count,
			DisplayedElements: len(processed),
		},
	})
}
