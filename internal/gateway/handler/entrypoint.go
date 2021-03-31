package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0tf/server/internal/shared"
	"strings"
	"time"
)

// Entrypoint handles all incoming gateway requests and delegates them to the corresponding element-specific handler
func Entrypoint(ctx *fiber.Ctx) error {
	// Extract required services
	namespaces := ctx.Locals("__services_namespaces").(shared.NamespaceService)
	elements := ctx.Locals("__services_elements").(shared.ElementService)

	// Retrieve the requested namespace
	namespaceID := strings.ToLower(ctx.Params("namespace_id"))
	namespace, err := namespaces.Namespace(namespaceID)
	if err != nil {
		return err
	}
	if namespace == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("the requested namespace does not exist")
	}

	// Retrieve the requested element
	elementKey := strings.ToLower(ctx.Params("element_key", "@"))
	element, err := elements.Element(namespace.ID, elementKey)
	if err != nil {
		return err
	}
	if element == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("the request element does not exist")
	}

	// Check if the requested element is already available
	if element.ValidFrom != -1 && element.ValidFrom > time.Now().Unix() {
		return ctx.Status(fiber.StatusPreconditionFailed).SendString("the requested element is not yet available")
	}

	// Check if the requested element is yet available
	if element.ValidUntil != -1 && element.ValidUntil < time.Now().Unix() {
		return ctx.Status(fiber.StatusPreconditionFailed).SendString("the requested element is not available anymore")
	}

	// Check if the requested element has already its maximum amount of views
	if element.MaxViews != -1 && element.Views >= element.MaxViews {
		return ctx.Status(fiber.StatusPreconditionFailed).SendString("the requested element cannot be viewed anymore")
	}

	// Increment the amount of views
	element.Views++
	if err := elements.CreateOrReplace(element); err != nil {
		return err
	}

	// Inject the element and delegate the request
	ctx.Locals("_element", element)
	switch element.Type {
	case shared.ElementTypePaste:
		return pasteHandler(ctx)
	case shared.ElementTypeRedirect:
		return redirectHandler(ctx)
	default:
		return ctx.Status(fiber.StatusNotImplemented).SendString("the type of the requested element is not yet supported")
	}
}
