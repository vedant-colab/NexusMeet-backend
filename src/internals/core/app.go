package core

import (
	"src/internals/routes"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupApp() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:               false,
		ServerHeader:          "Fiber/HighPerf",
		ReadTimeout:           5 * time.Second,
		WriteTimeout:          10 * time.Second,
		IdleTimeout:           15 * time.Second,
		ReadBufferSize:        8192,
		WriteBufferSize:       8192,
		DisableStartupMessage: true,
	})

	// Setup routes
	routes.SetupRoutes(app)

	return app
}
