package router

import (
	"davigo/s3uploader/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", handlers.HandleHealthCheck)
	app.Post("/", handlers.HandleUpload)
}
