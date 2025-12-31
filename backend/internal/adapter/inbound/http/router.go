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
	rateLimiter fiber.Handler,
) {
	// Item routes
	itemRoutes := app.Group("/api/items")
	itemRoutes.Get("/", itemHandler.GetAllItems)
	itemRoutes.Get("/:id", itemHandler.GetItemByID)
	itemRoutes.Post("/",
		middleware.AuthMiddleware(),
		rateLimiter,
		middleware.AdminOnly(),
		itemHandler.Create)
	itemRoutes.Put("/",
		middleware.AuthMiddleware(),
		rateLimiter,
		middleware.AdminOnly(),
		itemHandler.PutUpdate)
	itemRoutes.Patch("/:id",
		middleware.AuthMiddleware(),
		rateLimiter,
		middleware.AdminOnly(),
		itemHandler.PatchUpdate)
	itemRoutes.Delete("/:id",
		middleware.AuthMiddleware(),
		rateLimiter,
		middleware.AdminOnly(),
		itemHandler.Delete)

	// User routes
	userRoutes := app.Group("/api/users")
	userRoutes.Put("/:id",
		middleware.AuthMiddleware(),
		rateLimiter,
		middleware.AdminOrOwnerOnly(),
		userHandler.Update)

	userRoutes.Delete("/:id",
		middleware.AuthMiddleware(),
		rateLimiter,
		middleware.AdminOnly(),
		userHandler.Delete)
	userRoutes.Get("/",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
		userHandler.GetAllUsers)
	userRoutes.Get("/me",
		middleware.AuthMiddleware(),
		rateLimiter,
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
		rateLimiter,
		borrowingHandler.BorrowItem)
	borrowingRoutes.Post("/approve",
		middleware.AuthMiddleware(),
		rateLimiter,
		middleware.AdminOnly(),
		borrowingHandler.ApproveBorrowing)
	borrowingRoutes.Post("/reject",
		middleware.AuthMiddleware(),
		rateLimiter,
		middleware.AdminOnly(),
		borrowingHandler.RejectBorrowing)
	borrowingRoutes.Get("/status/:status",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
		borrowingHandler.GetBorrowingsByBorrowingStatus)
	borrowingRoutes.Get("/approval-status/:status",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
		borrowingHandler.GetBorrowingsByApprovalStatus)
	borrowingRoutes.Get("/user/:user_id",
		middleware.AuthMiddleware(),
		borrowingHandler.GetBorrowingByUserID)
	borrowingRoutes.Get("/stats",
		middleware.AuthMiddleware(),
		rateLimiter,
		borrowingHandler.UserStats)
}
