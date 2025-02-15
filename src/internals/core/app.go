package core

import (
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupApp() *fiber.App {
	runtime.GOMAXPROCS(4)
	app := fiber.New(fiber.Config{
		Prefork:               false,
		ServerHeader:          "Fiber/HighPerf",
		ReadTimeout:           5 * time.Second,
		WriteTimeout:          10 * time.Second,
		IdleTimeout:           15 * time.Second,
		ReadBufferSize:        8192,
		WriteBufferSize:       8192,
		DisableStartupMessage: false,
	})

	return app
}
