package gateway

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0tf/server/internal/shared"
	"strings"
)

// baseHandler looks up the element and delegates the request to the corresponding type handler
func baseHandler(ctx *fiber.Ctx) error {
	// Retrieve the namespace
	namespaces := ctx.Locals("__namespaces").(shared.NamespaceService)
	namespaceID := strings.ToLower(ctx.Params("namespace"))
	namespace, err := namespaces.Namespace(namespaceID)
	if err != nil {
		return err
	}
	if namespace == nil {
		return fiber.NewError(fiber.StatusNotFound, "the requested namespace does not exist")
	}

	// Retrieve the element
	elements := ctx.Locals("__elements").(shared.ElementService)
	elementKey := strings.ToLower(ctx.Params("key"))
	element, err := elements.Element(namespace.ID, elementKey)
	if err != nil {
		return err
	}
	if element == nil {
		return fiber.NewError(fiber.StatusNotFound, "the requested element does not exist")
	}

	// Inject the element and delegate the request
	ctx.Locals("_element", element)
	switch element.Type {
	case shared.ElementTypePaste:
		return pasteHandler(ctx)
	case shared.ElementTypeRedirect:
		return redirectHandler(ctx)
	default:
		return fiber.NewError(fiber.StatusNotImplemented, "unimplemented element type")
	}
}

// pasteHandler handles paste elements
func pasteHandler(ctx *fiber.Ctx) error {
	return ctx.SendString(ctx.Locals("_element").(*shared.Element).Data)
}

// redirectHandler handles paste elements
func redirectHandler(ctx *fiber.Ctx) error {
	return ctx.Redirect(ctx.Locals("_element").(*shared.Element).Data, fiber.StatusTemporaryRedirect)
}
