package gateway

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	recov "github.com/gofiber/fiber/v2/middleware/recover"
	log "github.com/sirupsen/logrus"
	"github.com/x0tf/server/internal/gateway/handler"
	"github.com/x0tf/server/internal/shared"
)

// Gateway represents the element-exposing gateway
type Gateway struct {
	app      *fiber.App
	Settings *Settings
	Services *Services
}

// Settings contains all settings important for the gateway
type Settings struct {
	Address      string
	Production   bool
	RootRedirect string
}

// Services contains all services used by the gateway
type Services struct {
	Namespaces shared.NamespaceService
	Elements   shared.ElementService
}

// Serve serves the gateway
func (gateway *Gateway) Serve() error {
	app := fiber.New(fiber.Config{
		DisableKeepalive:      true,
		DisableStartupMessage: gateway.Settings.Production,
	})

	// Enable panic recovering
	app.Use(recov.New())

	// Inject debug middlewares if the application runs in development mode
	if !gateway.Settings.Production {
		app.Use(logger.New())
		app.Use(pprof.New())
	}

	// Inject the application data
	app.Use(func(ctx *fiber.Ctx) error {
		ctx.Locals("__services_namespaces", gateway.Services.Namespaces)
		ctx.Locals("__services_elements", gateway.Services.Elements)

		return ctx.Next()
	})

	app.Get("/:namespace_id/:element_key?", handler.Entrypoint)

	// Define the root redirect
	if gateway.Settings.RootRedirect != "" {
		app.Get("/", func(ctx *fiber.Ctx) error {
			return ctx.Redirect(gateway.Settings.RootRedirect, fiber.StatusPermanentRedirect)
		})
	}

	log.WithField("address", gateway.Settings.Address).Info("Serving the gateway")
	gateway.app = app
	return app.Listen(gateway.Settings.Address)
}

// Shutdown gracefully shuts down the gateway
func (gateway *Gateway) Shutdown() error {
	log.Info("Shutting down the gateway")
	return gateway.app.Shutdown()
}
