package gateway

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	recov "github.com/gofiber/fiber/v2/middleware/recover"
	log "github.com/sirupsen/logrus"
	"github.com/x0tf/server/internal/shared"
)

// Gateway represents the element-exposing gateway
type Gateway struct {
	app          *fiber.App
	Address      string
	Production   bool
	Namespaces   shared.NamespaceService
	Elements     shared.ElementService
	RootRedirect string
}

// Serve serves the gateway
func (gateway *Gateway) Serve() error {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: gateway.Production,
	})

	// Enable panic recovering
	app.Use(recov.New())

	// Inject debug middlewares if the application runs in development mode
	if !gateway.Production {
		app.Use(logger.New())
		app.Use(pprof.New())
	}

	// Inject the application data
	app.Use(func(ctx *fiber.Ctx) error {
		ctx.Locals("__namespaces", gateway.Namespaces)
		ctx.Locals("__elements", gateway.Elements)
		return ctx.Next()
	})

	app.Get("/:namespace/:key?", baseHandler)

	// Define the root redirect
	if gateway.RootRedirect != "" {
		app.Get("/", func(ctx *fiber.Ctx) error {
			return ctx.Redirect(gateway.RootRedirect, fiber.StatusPermanentRedirect)
		})
	}

	log.WithField("address", gateway.Address).Info("Serving the gateway")
	gateway.app = app
	return app.Listen(gateway.Address)
}

// Shutdown gracefully shuts down the gateway
func (gateway *Gateway) Shutdown() error {
	log.Info("Shutting down the gateway")
	return gateway.app.Shutdown()
}
