package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0tf/server/internal/validation"
)

const (
	errorCodeGenericUnexpectedError   = 1000
	errorCodeGenericUnauthorized      = 1001
	errorCodeGenericBadQueryParameter = 1002
	errorCodeGenericBadRequestBody    = 1003
	errorCodeGenericNamespaceNotFound = 1004
	errorCodeGenericElementNotFound   = 1005

	errorCodeNamespaceIllegalNamespaceID = 2000
	errorCodeNamespaceNamespaceIDInUse   = 2001

	errorCodeElementElementKeyInUse      = 3000
	errorCodeElementNamespaceDeactivated = 3001

	errorCodeElementPasteEmptyPasteContent = 3100

	errorCodeElementRedirectInvalidTargetURL = 3200
)

var (
	errorGenericUnexpectedError = func(err error) *apiError {
		return newError(fiber.StatusInternalServerError, errorCodeGenericUnexpectedError, err.Error(), nil)
	}

	errorGenericUnauthorized = newError(fiber.StatusUnauthorized, errorCodeGenericUnauthorized, "unauthorized", nil)

	errorGenericBadQueryParameter = func(name, given, desiredType string) *apiError {
		return newError(fiber.StatusBadRequest, errorCodeGenericBadQueryParameter, "bad query parameter", fiber.Map{
			"name":         name,
			"given":        given,
			"desired_type": desiredType,
		})
	}

	errorGenericBadRequestBody = newError(fiber.StatusBadRequest, errorCodeGenericBadRequestBody, "bad request body", nil)

	errorGenericNamespaceNotFound = newError(fiber.StatusNotFound, errorCodeGenericNamespaceNotFound, "namespace not found", nil)

	errorGenericElementNotFound = newError(fiber.StatusNotFound, errorCodeGenericElementNotFound, "element not found", nil)

	errorNamespaceIllegalNamespaceID = func(violations []validation.NamespaceIDViolation) *apiError {
		return newError(fiber.StatusUnprocessableEntity, errorCodeNamespaceIllegalNamespaceID, "illegal namespace ID", fiber.Map{
			"violations": violations,
		})
	}

	errorNamespaceNamespaceIDInUse = newError(fiber.StatusConflict, errorCodeNamespaceNamespaceIDInUse, "namespace ID in use", nil)

	errorElementElementKeyInUse = newError(fiber.StatusConflict, errorCodeElementElementKeyInUse, "element key in use", nil)

	errorElementNamespaceDeactivated = newError(fiber.StatusForbidden, errorCodeElementNamespaceDeactivated, "this namespace is deactivated", nil)

	errorElementPasteEmptyPasteContent = newError(fiber.StatusUnprocessableEntity, errorCodeElementPasteEmptyPasteContent, "empty paste content", nil)

	errorElementRedirectInvalidTargetURL = newError(fiber.StatusUnprocessableEntity, errorCodeElementRedirectInvalidTargetURL, "invalid target URL", nil)
)
