package main

import (
	"log"
	"src/internals/config"
	"src/internals/core"
	"src/internals/database"
	"src/internals/routes"
	"time"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func main() {
	// Initialize the Fiber app
	config.LoadEnv()
	database.ConnectDB()
	app := core.SetupApp()

	server := app.Server()
	server.MaxConnsPerIP = 10000
	server.MaxIdleWorkerDuration = 30 * time.Second
	server.TCPKeepalive = true
	server.TCPKeepalivePeriod = 1 * time.Minute
	server.Concurrency = 256 * 1024

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, http://localhost:5174", // Your frontend origin
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS", // Allow necessary methods
		AllowCredentials: true,                              // Enable if using cookies or session data
	}))

	app.Use(helmet.New())

	app.Use(logger.New(logger.Config{
		Format:     "${time} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "Asia/Kolkata",
	}))

	app.Use(idempotency.New())

	app.Use(limiter.New(limiter.Config{
		Max:               20,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	routes.SetupRoutes(app)

	app.Get("/metrics", monitor.New(monitor.Config{
		Title: "NexusMeet metrics page",
	}))

	// Start the server
	log.Fatal(app.Listen(":8000"))
}
