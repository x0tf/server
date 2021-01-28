package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	recov "github.com/gofiber/fiber/v2/middleware/recover"
	v1 "github.com/x0tf/server/internal/api/v1"
	"github.com/x0tf/server/internal/shared"
)

// API represents the REST API
type API struct {
	Address    string
	Production bool
	Version    string
	Namespaces shared.NamespaceService
	Elements   shared.ElementService
}

// Serve serves the REST API
func (api *API) Serve() error {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: api.Production,
	})

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

	// Enable panic recovering
	app.Use(recov.New())

	// Inject debug middlewares if the application runs in development mode
	if !api.Production {
		app.Use(logger.New())
		app.Use(pprof.New())
	}

	// Inject the application data
	app.Use(func(ctx *fiber.Ctx) error {
		ctx.Locals("__production", api.Production)
		ctx.Locals("__version", api.Version)
		ctx.Locals("__namespaces", api.Namespaces)
		ctx.Locals("__elements", api.Elements)
		return ctx.Next()
	})

	// Route the v1 API endpoints
	v1router := app.Group("/v1")
	v1router.Get("/info", v1.EndpointGetInfo)

	return app.Listen(api.Address)
}
