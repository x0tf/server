package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0tf/server/internal/shared"
	"github.com/x0tf/server/internal/token"
	"github.com/x0tf/server/internal/utils"
	"github.com/x0tf/server/internal/validation"
)

// EndpointListNamespaces handles the GET /v1/namespaces endpoint
func EndpointListNamespaces(ctx *fiber.Ctx) error {
	namespaces := ctx.Locals("__namespaces").(shared.NamespaceService)
	list, err := namespaces.Namespaces()
	if err != nil {
		return err
	}

	// Process the list to remove the tokens
	processedList := make([]shared.Namespace, 0, len(list))
	for _, namespace := range list {
		processedNamespace := *namespace
		processedNamespace.Token = ""
		processedList = append(processedList, processedNamespace)
	}
	return ctx.JSON(processedList)
}

// EndpointGetNamespace handles the GET /v1/namespaces/:namespace endpoint
func EndpointGetNamespace(ctx *fiber.Ctx) error {
	namespace := *(ctx.Locals("_namespace").(*shared.Namespace))
	namespace.Token = ""
	return ctx.JSON(namespace)
}

// EndpointCreateNamespace handles the POST /v1/namespaces/:namespace endpoint
func EndpointCreateNamespace(ctx *fiber.Ctx) error {
	// Check if the user has to provide an invite code
	invites, _ := ctx.Locals("__invites").(shared.InviteService)
	var usedInvite string
	if invites != nil && !ctx.Locals("_admin").(bool) {
		// Parse the JSON body into a map
		var data map[string]interface{}
		if err := ctx.BodyParser(&data); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "could not parse request body")
		}

		// Try to retrieve the invite code value of the JSON body
		invite, ok := data["invite"].(string)
		if !ok {
			return fiber.NewError(fiber.StatusBadRequest, "got an illegal value as invite code")
		}

		// Check if the given invite code is valid
		isValid, err := invites.IsValid(invite)
		if err != nil {
			return err
		}
		if !isValid {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid invite code")
		}
		usedInvite = invite
	}

	// Validate the given namespace ID
	id := ctx.Params("namespace")
	if errors := validation.ValidateNamespaceID(id); len(errors) > 0 {
		messages := make([]string, 0, len(errors))
		for _, err := range errors {
			messages = append(messages, err.Error())
		}
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"messages": messages,
		})
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
		Token:  utils.GenerateToken(),
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

	// Delete the invite code if one was used
	if usedInvite != "" && invites != nil {
		if err = invites.Delete(usedInvite); err != nil {
			return err
		}
	}

	// Return the copied namespace with the raw token still placed in it
	return ctx.JSON(namespaceCopy)
}

// EndpointResetNamespaceToken handles the POST /v1/namespaces/:namespace/resetToken endpoint
func EndpointResetNamespaceToken(ctx *fiber.Ctx) error {
	namespace := ctx.Locals("_namespace").(*shared.Namespace)
	newToken := utils.GenerateToken()
	hash, err := token.Hash(newToken)
	if err != nil {
		return err
	}
	namespace.Token = hash

	namespaces := ctx.Locals("__namespaces").(shared.NamespaceService)
	if err = namespaces.CreateOrReplace(namespace); err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{"token": newToken})
}

// EndpointDeactivateNamespace handles the POST /v1/namespaces/:namespace/deactivate endpoint
func EndpointDeactivateNamespace(ctx *fiber.Ctx) error {
	namespace := ctx.Locals("_namespace").(*shared.Namespace)
	namespaces := ctx.Locals("__namespaces").(shared.NamespaceService)
	namespace.Active = false
	return namespaces.CreateOrReplace(namespace)
}

// EndpointActivateNamespace handles the POST /v1/namespaces/:namespace/activate endpoint
func EndpointActivateNamespace(ctx *fiber.Ctx) error {
	namespace := ctx.Locals("_namespace").(*shared.Namespace)
	namespaces := ctx.Locals("__namespaces").(shared.NamespaceService)
	namespace.Active = true
	return namespaces.CreateOrReplace(namespace)
}

// EndpointDeleteNamespace handles the DELETE /v1/namespaces/:namespace endpoint
func EndpointDeleteNamespace(ctx *fiber.Ctx) error {
	namespace := ctx.Locals("_namespace").(*shared.Namespace)
	namespaces := ctx.Locals("__namespaces").(shared.NamespaceService)
	elements := ctx.Locals("__elements").(shared.ElementService)

	if err := elements.DeleteInNamespace(namespace.ID); err != nil {
		return err
	}

	return namespaces.Delete(namespace.ID)
}
