package main

import (
	"log"
	"src/internals/config"
	"src/internals/core"
	"src/internals/database"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Initialize the Fiber app
	config.LoadEnv()
	database.ConnectDB()
	app := core.SetupApp()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Start the server
	log.Fatal(app.Listen(":8000"))
}
