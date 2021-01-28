package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/x0tf/server/internal/shared"
)

// API represents the REST API
type API struct {
	Address    string
	Namespaces shared.NamespaceService
	Elements   shared.ElementService
}

// Serve serves the REST API
func (api *API) Serve() error {
	app := fiber.New()

	// Include CORS response headers
	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: true,
		ExposeHeaders:    "",
		MaxAge:           0,
	}))

	// Add debug route
	// TODO: Remove this
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("debug route")
	})

	return app.Listen(api.Address)
}
