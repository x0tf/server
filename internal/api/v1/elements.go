package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0tf/server/internal/shared"
	"github.com/x0tf/server/internal/utils"
	"net/url"
	"strings"
)

// EndpointListElements handles the GET /v1/elements endpoint
func EndpointListElements(ctx *fiber.Ctx) error {
	elements := ctx.Locals("__elements").(shared.ElementService)
	list, err := elements.Elements()
	if err != nil {
		return err
	}
	if list == nil {
		list = []*shared.Element{}
	}
	return ctx.JSON(list)
}

// EndpointListNamespaceElements handles the GET /v1/elements/:namespace endpoint
func EndpointListNamespaceElements(ctx *fiber.Ctx) error {
	namespace := ctx.Locals("_namespace").(*shared.Namespace)
	elements := ctx.Locals("__elements").(shared.ElementService)
	list, err := elements.ElementsInNamespace(namespace.ID)
	if err != nil {
		return err
	}
	if list == nil {
		list = []*shared.Element{}
	}
	return ctx.JSON(list)
}

// EndpointGetElement handles the GET /v1/elements/:namespace/:key endpoint
func EndpointGetElement(ctx *fiber.Ctx) error {
	namespace := ctx.Locals("_namespace").(*shared.Namespace)
	elements := ctx.Locals("__elements").(shared.ElementService)
	element, err := elements.Element(namespace.ID, strings.ToLower(ctx.Params("key")))
	if err != nil {
		return err
	}
	if element == nil {
		return fiber.NewError(fiber.StatusNotFound, "that element does not exist")
	}
	return ctx.JSON(element)
}

// EndpointCreatePasteElement handles the POST /v1/elements/:namespace/paste/:key? endpoint
func EndpointCreatePasteElement(ctx *fiber.Ctx) error {
	isAdmin := ctx.Locals("_admin").(bool)
	namespace := ctx.Locals("_namespace").(*shared.Namespace)
	elements := ctx.Locals("__elements").(shared.ElementService)

	// Check if the namespace is deactivated
	if !namespace.Active && !isAdmin {
		return fiber.NewError(fiber.StatusForbidden, "this namespace is deactivated")
	}

	// Parse the JSON body into a map
	var data map[string]interface{}
	if err := ctx.BodyParser(&data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "could not parse request body")
	}

	// Read and validate the paste content out of the request body
	content, ok := data["content"].(string)
	if !ok || strings.TrimSpace(content) == "" {
		return fiber.NewError(fiber.StatusBadRequest, "got an illegal or empty value as paste content")
	}

	// Generate a new element key
	key := strings.TrimSpace(strings.ToLower(ctx.Params("key")))
	if key == "" {
		for {
			key = utils.GenerateElementKey()
			found, err := elements.Element(namespace.ID, key)
			if err != nil {
				return err
			}
			if found == nil {
				break
			}
		}
	}
	found, err := elements.Element(namespace.ID, key)
	if err != nil {
		return err
	}
	if found != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "the given element key is already in use")
	}

	// Create the element
	element := &shared.Element{
		Namespace: namespace.ID,
		Key:       key,
		Type:      shared.ElementTypePaste,
		Data:      content,
	}
	if err = elements.CreateOrReplace(element); err != nil {
		return err
	}
	return ctx.JSON(element)
}

// EndpointCreateRedirectElement handles the POST /v1/elements/:namespace/redirect/:key? endpoint
func EndpointCreateRedirectElement(ctx *fiber.Ctx) error {
	isAdmin := ctx.Locals("_admin").(bool)
	namespace := ctx.Locals("_namespace").(*shared.Namespace)
	elements := ctx.Locals("__elements").(shared.ElementService)

	// Check if the namespace is deactivated
	if !namespace.Active && !isAdmin {
		return fiber.NewError(fiber.StatusForbidden, "this namespace is deactivated")
	}

	// Parse the JSON body into a map
	var data map[string]interface{}
	if err := ctx.BodyParser(&data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "could not parse request body")
	}

	// Read and validate the target URL out of the request body
	target, ok := data["target"].(string)
	if !ok || strings.TrimSpace(target) == "" {
		return fiber.NewError(fiber.StatusBadRequest, "got an illegal or empty value as target URL")
	}
	parsedURL, err := url.Parse(target)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") || parsedURL.Host == "" {
		return fiber.NewError(fiber.StatusBadRequest, "the given target URL is no http(s) url")
	}

	// Generate a new element key
	key := strings.TrimSpace(strings.ToLower(ctx.Params("key")))
	if key == "" {
		for {
			key = utils.GenerateElementKey()
			found, err := elements.Element(namespace.ID, key)
			if err != nil {
				return err
			}
			if found == nil {
				break
			}
		}
	}
	found, err := elements.Element(namespace.ID, key)
	if err != nil {
		return err
	}
	if found != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "the given element key is already in use")
	}

	// Create the element
	element := &shared.Element{
		Namespace: namespace.ID,
		Key:       key,
		Type:      shared.ElementTypeRedirect,
		Data:      parsedURL.String(),
	}
	if err = elements.CreateOrReplace(element); err != nil {
		return err
	}
	return ctx.JSON(element)
}

// EndpointDeleteElement handles the DELETE /v1/elements/:namespace/:key endpoint
func EndpointDeleteElement(ctx *fiber.Ctx) error {
	namespace := ctx.Locals("_namespace").(*shared.Namespace)
	elements := ctx.Locals("__elements").(shared.ElementService)
	element, err := elements.Element(namespace.ID, strings.ToLower(ctx.Params("key")))
	if err != nil {
		return err
	}
	if element == nil {
		return fiber.NewError(fiber.StatusNotFound, "that element does not exist")
	}
	return elements.Delete(namespace.ID, element.Key)
}
