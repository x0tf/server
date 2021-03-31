package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0tf/server/internal/shared"
	"github.com/x0tf/server/internal/utils"
	"strings"
	"time"
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

// EndpointGetElement handles the 'GET /v2/elements/:namespace_id/:element_id' endpoint
func EndpointGetElement(ctx *fiber.Ctx) error {
	element := *(ctx.Locals("_element").(*shared.Element))
	element.InternalData = map[string]interface{}{}
	return ctx.JSON(element)
}

type endpointCreatePasteElementRequestBody struct {
	Key        *string `json:"key" xml:"key" form:"key"`
	MaxViews   *int    `json:"max_views" xml:"max_views" form:"max_views"`
	ValidFrom  *int64  `json:"valid_from" xml:"valid_from" form:"valid_from"`
	ValidUntil *int64  `json:"valid_until" xml:"valid_until" form:"valid_until"`
	Content    string  `json:"content" xml:"content" form:"content"`
}

// EndpointCreatePasteElement handles the 'POST /v2/elements/:namespace_id/paste' endpoint
func EndpointCreatePasteElement(ctx *fiber.Ctx) error {
	// Extract required services
	elements := ctx.Locals("__services_elements").(shared.ElementService)

	// Extract required resources
	namespace := ctx.Locals("_namespace").(*shared.Namespace)

	// Try to parse the body into a request body struct
	body := new(endpointCreatePasteElementRequestBody)
	if err := ctx.BodyParser(body); err != nil {
		return newError(fiber.StatusBadRequest, errorCodeGenericBadRequestBody, "invalid request body", nil)
	}

	// Validate the content length
	if len(body.Content) == 0 {
		return newError(fiber.StatusBadRequest, errorCodeElementEmptyPasteContent, "empty paste content", nil)
	}

	// Validate the wished key if one was provided or generate a new one
	var key string
	if body.Key != nil {
		wishedKey := strings.ToLower(strings.TrimSpace(*body.Key))
		found, err := elements.Element(namespace.ID, wishedKey)
		if err != nil {
			return err
		}
		if found != nil {
			return newError(fiber.StatusConflict, errorCodeElementElementKeyInUse, "element key in use", nil)
		}
		key = wishedKey
	} else {
		for {
			generated := utils.GenerateElementKey()
			found, err := elements.Element(namespace.ID, generated)
			if err != nil {
				return err
			}
			if found == nil {
				key = generated
				break
			}
		}
	}

	// Define the maximum amount of views
	maxViews := -1
	if body.MaxViews != nil {
		maxViews = *body.MaxViews
	}

	// Define the timestamp when the element should become valid
	validFrom := int64(-1)
	if body.ValidFrom != nil {
		validFrom = *body.ValidFrom
	}

	// Define the timestamp when the element should expire
	validUntil := int64(-1)
	if body.ValidUntil != nil {
		validUntil = *body.ValidUntil
	}

	// Create and respond with the element
	element := &shared.Element{
		Namespace:    namespace.ID,
		Key:          key,
		Type:         shared.ElementTypePaste,
		InternalData: map[string]interface{}{},
		PublicData: map[string]interface{}{
			"content": body.Content,
		},
		Views:      0,
		MaxViews:   maxViews,
		ValidFrom:  validFrom,
		ValidUntil: validUntil,
		Created:    time.Now().Unix(),
	}
	if err := elements.CreateOrReplace(element); err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(element)
}