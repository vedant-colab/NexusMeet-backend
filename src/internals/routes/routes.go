package routes

import (
	"src/internals/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// WebSocket route
	app.Get("/ws/:room_id/:user_id", controllers.HandleWebSocket)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	users := v1.Group("/users")
	users.Post("/signup", controllers.Signup)
	users.Post("/signin", controllers.Signin)
	users.Get("/logout", controllers.Logout)
	users.Get("/heavy1", controllers.Load)
	users.Get("/users", controllers.GetUsers)

	profiles := v1.Group("/profiles")
	profiles.Get("/:username", controllers.FetchProfile)
	profiles.Put("/:username", controllers.UpdateProfile)
}
