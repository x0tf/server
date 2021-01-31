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
	Address     string
	Production  bool
	Version     string
	AdminTokens []string
	Namespaces  shared.NamespaceService
	Elements    shared.ElementService
	Invites     shared.InviteService
}

// Serve serves the REST API
func (api *API) Serve() error {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: api.Production,
		ErrorHandler:          errorHandler,
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
		if api.Invites != nil {
			ctx.Locals("__invites", api.Invites)
		}
		ctx.Locals("__admin_tokens", api.AdminTokens)
		return ctx.Next()
	})

	// Route the v1 API endpoints
	v1router := app.Group("/v1")
	{
		v1router.Get("/info", v1.EndpointGetInfo)

		// Register the invite endpoints if required
		if api.Invites != nil {
			v1router.Post("/invites/:code?", v1.MiddlewareAdminAuth, v1.MiddlewareRequireAdminAuth, v1.EndpointCreateInvite)
		}

		// Register the namespace endpoints
		v1router.Get("/namespaces/:namespace", v1.MiddlewareInjectNamespace, v1.EndpointGetNamespace)
		v1router.Post("/namespaces/:namespace", v1.EndpointCreateNamespace)
		v1router.Post("/namespaces/:namespace/resetToken", v1.MiddlewareAdminAuth, v1.MiddlewareInjectNamespace, v1.MiddlewareTokenAuth, v1.EndpointResetNamespaceToken)
		v1router.Post("/namespaces/:namespace/deactivate", v1.MiddlewareAdminAuth, v1.MiddlewareRequireAdminAuth, v1.MiddlewareInjectNamespace, v1.EndpointDeactivateNamespace)
		v1router.Post("/namespaces/:namespace/activate", v1.MiddlewareAdminAuth, v1.MiddlewareRequireAdminAuth, v1.MiddlewareInjectNamespace, v1.EndpointActivateNamespace)
		v1router.Delete("/namespaces/:namespace", v1.MiddlewareAdminAuth, v1.MiddlewareInjectNamespace, v1.MiddlewareTokenAuth, v1.EndpointDeleteNamespace)
	}

	return app.Listen(api.Address)
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if fiberError, ok := err.(*fiber.Error); ok {
		code = fiberError.Code
	}
	return ctx.Status(code).JSON(fiber.Map{
		"messages": []string{err.Error()},
	})
}
