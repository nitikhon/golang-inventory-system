package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nitikhon/golang-inventory-system/internal/adapter/inbound/http/middleware"
)

func SetupRoutes(
	app *fiber.App,
	itemHandler *ItemHandler,
	userHandler *UserHandler,
	borrowingHandler *BorrowingHandler,
) {
	// Item routes
	itemRoutes := app.Group("/api/items")
	itemRoutes.Get("/", itemHandler.GetAllItems)
	itemRoutes.Get("/:id", itemHandler.GetItemByID)
	itemRoutes.Post("/",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
		itemHandler.Create)
	itemRoutes.Put("/", 
		middleware.AuthMiddleware(), 
		middleware.AdminOnly(), 
		itemHandler.Update)
	itemRoutes.Delete("/:id",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
		itemHandler.Delete)

	// User routes
	userRoutes := app.Group("/api/users")
	userRoutes.Put("/",
		middleware.AuthMiddleware(),
		middleware.AdminOrOwnerOnly(),
		userHandler.Update)

	userRoutes.Delete("/:id",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
		userHandler.Delete)
	userRoutes.Get("/",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
		userHandler.GetAllUsers)
	userRoutes.Get("/me",
		middleware.AuthMiddleware(),
		userHandler.Me) // define this before :id to prevent an error
	userRoutes.Get("/:id",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
		userHandler.GetUserByID)
	userRoutes.Get("/username/:username",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
		userHandler.GetUserByUsername)
	userRoutes.Get("/email/:email",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
		userHandler.GetUserByEmail)
	userRoutes.Get("/phone/:phone",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
		userHandler.GetUserByPhone)

	// Auth routes
	userRoutes.Post("/login", userHandler.Login)
	userRoutes.Post("/register", userHandler.Create)
	userRoutes.Post("/refresh-token", userHandler.RefreshToken)
	userRoutes.Post("/logout",
		middleware.AuthMiddleware(),
		userHandler.Logout)

	// Borrowing routes
	borrowingRoutes := app.Group("/api/borrows")
	borrowingRoutes.Post("/",
		middleware.AuthMiddleware(),
		borrowingHandler.BorrowItem)
	borrowingRoutes.Post("/approve",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
		borrowingHandler.ApproveBorrowing)
	borrowingRoutes.Post("/reject",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
		borrowingHandler.RejectBorrowing)
}
