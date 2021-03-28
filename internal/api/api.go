package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	recov "github.com/gofiber/fiber/v2/middleware/recover"
	log "github.com/sirupsen/logrus"
	v2 "github.com/x0tf/server/internal/api/v2"
	"github.com/x0tf/server/internal/shared"
)

// API represents the REST API
type API struct {
	app      *fiber.App
	Settings *Settings
	Services *Services
}

// Settings contains all settings important for the REST API
type Settings struct {
	Address           string
	RequestsPerMinute int
	Production        bool
	Version           string
	AdminTokens       []string
	InvitesEnabled    bool
}

// Services contains all services used by the REST API
type Services struct {
	Namespaces shared.NamespaceService
	Elements   shared.ElementService
	Invites    shared.InviteService
}

// Serve serves the REST API
func (api *API) Serve() error {
	app := fiber.New(fiber.Config{
		ErrorHandler:          v2.ErrorHandler,
		DisableKeepalive:      true,
		DisableStartupMessage: api.Settings.Production,
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
	if !api.Settings.Production {
		app.Use(logger.New())
		app.Use(pprof.New())
	}

	// Inject the rate limiter middleware
	app.Use(limiter.New(limiter.Config{
		Next: func(_ *fiber.Ctx) bool {
			return !api.Settings.Production
		},
		Max: api.Settings.RequestsPerMinute,
		LimitReached: func(ctx *fiber.Ctx) error {
			return fiber.ErrTooManyRequests
		},
	}))

	// Inject the application data
	app.Use(func(ctx *fiber.Ctx) error {
		ctx.Locals("__settings_address", api.Settings.Address)
		ctx.Locals("__settings_requests_per_minute", api.Settings.RequestsPerMinute)
		ctx.Locals("__settings_production", api.Settings.Production)
		ctx.Locals("__settings_version", api.Settings.Version)
		ctx.Locals("__settings_admin_tokens", api.Settings.AdminTokens)
		ctx.Locals("__settings_invites_enabled", api.Settings.InvitesEnabled)

		ctx.Locals("__services_namespaces", api.Services.Namespaces)
		ctx.Locals("__services_elements", api.Services.Elements)
		ctx.Locals("__services_invites", api.Services.Invites)

		return ctx.Next()
	})

	// Mark the v1 API as replaced
	app.All("/v1/*", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).SendString("API version v1 has been replaced by version v2")
	})

	// Route the v2 API endpoints
	v2group := app.Group("/v2")
	{
		v2group.Get("/info", v2.EndpointGetInfo)

		v2group.Get("/namespaces", v2.MiddlewareAdminAuth(true), v2.EndpointGetNamespaces)
	}

	log.WithField("address", api.Settings.Address).Info("Serving the REST API")
	api.app = app
	return app.Listen(api.Settings.Address)
}

// Shutdown gracefully shuts down the REST API
func (api *API) Shutdown() error {
	log.Info("Shutting down the REST API")
	return api.app.Shutdown()
}
