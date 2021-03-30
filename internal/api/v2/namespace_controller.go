package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0tf/server/internal/shared"
	"github.com/x0tf/server/internal/token"
	"github.com/x0tf/server/internal/utils"
	"github.com/x0tf/server/internal/validation"
	"time"
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

type endpointCreateNamespaceRequestBody struct {
	ID string `json:"id" xml:"id" form:"id"`
}

// EndpointCreateNamespace handles the 'POST /v2/namespaces' endpoint
func EndpointCreateNamespace(ctx *fiber.Ctx) error {
	// Extract required services
	namespaces := ctx.Locals("__services_namespaces").(shared.NamespaceService)

	// Try to parse the body into a request body struct
	body := new(endpointCreateNamespaceRequestBody)
	if err := ctx.BodyParser(body); err != nil {
		return newError(fiber.StatusBadRequest, errorCodeGenericBadRequestBody, "invalid request body", nil)
	}

	// Validate the given ID
	violations := validation.ValidateNamespaceID(body.ID)
	if len(violations) > 0 {
		return newError(fiber.StatusUnprocessableEntity, errorCodeNamespaceIllegalNamespaceID, "illegal namespace ID", fiber.Map{
			"violations": violations,
		})
	}

	// Check if a namespace with that ID already exists
	found, err := namespaces.Namespace(body.ID)
	if err != nil {
		return err
	}
	if found != nil {
		return newError(fiber.StatusConflict, errorCodeNamespaceNamespaceIDInUse, "namespace ID in use", nil)
	}

	// Create a new token for the namespace
	generatedToken := utils.GenerateToken()
	hashedToken, err := token.Hash(generatedToken)
	if err != nil {
		return err
	}

	// Create a namespace instance with default values and take a copy from it to include the raw token
	namespace := &shared.Namespace{
		ID:      body.ID,
		Token:   hashedToken,
		Active:  true,
		Created: time.Now().Unix(),
	}
	copy := *namespace
	copy.Token = generatedToken

	// Insert the created instance into the database and respond with the copy
	if err := namespaces.CreateOrReplace(namespace); err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(copy)
}

type endpointPatchNamespaceRequestBody struct {
	Active *bool `json:"active" xml:"active" form:"active"`
}

// EndpointPatchNamespace handles the 'PATCH /v2/namespaces/:namespace_id' endpoint
func EndpointPatchNamespace(ctx *fiber.Ctx) error {
	// Extract required services
	namespaces := ctx.Locals("__services_namespaces").(shared.NamespaceService)

	// Extract required resources
	namespace := ctx.Locals("_namespace").(*shared.Namespace)

	// Try to parse the body into a request body struct
	body := new(endpointPatchNamespaceRequestBody)
	if err := ctx.BodyParser(body); err != nil {
		return newError(fiber.StatusBadRequest, errorCodeGenericBadRequestBody, "invalid request body", nil)
	}

	// Update the namespace accordingly
	if body.Active != nil {
		namespace.Active = *body.Active
	}
	if err := namespaces.CreateOrReplace(namespace); err != nil {
		return err
	}
	copy := *namespace
	copy.Token = ""
	return ctx.JSON(copy)
}
