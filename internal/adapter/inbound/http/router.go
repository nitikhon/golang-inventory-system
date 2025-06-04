package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nitikhon/golang-inventory-system/internal/adapter/inbound/http/middleware"
)

func SetupRoutes(app *fiber.App, itemHandler *ItemHandler, userHandler *UserHandler) {
	// Item routes
	itemRoutes := app.Group("/items")
	itemRoutes.Get("/", itemHandler.GetAllItems)
	itemRoutes.Get("/:id", itemHandler.GetItemByID)
	itemRoutes.Post("/", itemHandler.Create)
	itemRoutes.Put("/", itemHandler.Update)
	itemRoutes.Delete("/:id", itemHandler.Delete)

	// User routes
	userRoutes := app.Group("/users")
	userRoutes.Put("/", userHandler.Update)
	
	userRoutes.Delete("/:id", userHandler.Delete)
	userRoutes.Get("/", userHandler.GetAllUsers)
	userRoutes.Get("/me", middleware.AuthMiddleware(), userHandler.Me) // define this before :id to prevent an error
	userRoutes.Get("/:id", userHandler.GetUserByID)
	userRoutes.Get("/username/:username", userHandler.GetUserByUsername)
	userRoutes.Get("/email/:email", userHandler.GetUserByEmail)
	userRoutes.Get("/phone/:phone", userHandler.GetUserByPhone)

	// Auth routes
	userRoutes.Post("/login", userHandler.Login)
	userRoutes.Post("/register", userHandler.Create) // Register user
	userRoutes.Post("/refresh-token", userHandler.RefreshToken)
	

	// Protected routes
	// protectedRoutes := app.Group("/protected", middleware.AuthMiddleware())
}
