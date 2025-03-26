package http

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, itemHandler *ItemHandler) {
	app.Get("/items", itemHandler.FindAll)
	app.Get("/items/:id", itemHandler.FindByID)
	app.Post("/items", itemHandler.Create)
	app.Put("/items", itemHandler.Update)
	app.Delete("/items/:id", itemHandler.Delete)
}