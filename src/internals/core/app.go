package core

import (
	"src/internals/routes"

	"github.com/gofiber/fiber/v2"
)

func SetupApp() *fiber.App {
	app := fiber.New()

	// Setup routes
	routes.SetupRoutes(app)

	return app
}
