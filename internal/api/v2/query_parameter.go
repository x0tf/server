package v2

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// parseQueryInt parses a query parameter into an integer
func parseQueryInt(ctx *fiber.Ctx, name string, defaultVale int) (int, error) {
	value := ctx.Query(name)
	if value == "" {
		return defaultVale, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, errorGenericBadQueryParameter(name, value, "int")
	}
	return parsed, nil
}
