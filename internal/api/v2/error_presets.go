package v2

import "github.com/gofiber/fiber/v2"

var (
	errUnauthorized = newError(fiber.StatusUnauthorized, errorCodeGenericUnauthorized, "unauthorized", nil)
	errForbidden    = newError(fiber.StatusForbidden, errorCodeGenericForbidden, "forbidden", nil)
)
