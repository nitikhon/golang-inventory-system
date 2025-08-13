package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/nitikhon/golang-inventory-system/internal/adapter/inbound/http"
	"github.com/nitikhon/golang-inventory-system/internal/adapter/outbound/database"
	"github.com/nitikhon/golang-inventory-system/internal/adapter/outbound/repository"
	"github.com/nitikhon/golang-inventory-system/internal/config"
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/core/service"
	"github.com/nitikhon/golang-inventory-system/internal/util/seed"
)

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	// Initialize configuration
	config := config.NewConfig()

	// Connect to the database
	db, err := database.NewDatabase(*config)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	// Drop all tables before migration
    err = db.Migrator().DropTable(&entity.Item{}, &entity.User{}, &entity.Borrowing{})
    if err != nil {
        log.Fatal("Error dropping tables: ", err)
    }
    log.Println("All tables dropped successfully.")

    // Auto-migrate database models
    err = db.AutoMigrate(&entity.Item{}, &entity.User{}, &entity.Borrowing{})
    if err != nil {
        log.Fatal("Error migrating tables: ", err)
    }
    log.Println("Database migrated successfully.")

	// Seed the database if it's empty
	if err := seed.SeedDB(db); err != nil {
		log.Println("Error: ", err)
	} else {
		log.Println("Database is empty, seeding...")
	}

	// Initialize repositories, services, and handlers
	itemRepo := repository.NewItemRepository(db)
	itemService := service.NewItemService(itemRepo)
	itemHandler := http.NewItemHandler(itemService)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := http.NewUserHandler(userService)

	borrowingRepo := repository.NewBorrowingRepository(db)
	borrowingService := service.NewBorrowingService(borrowingRepo, itemRepo, userRepo)
	borrowingHandler := http.NewBorrowingHandler(borrowingService)

	// Create a new Fiber app
	app := fiber.New()

	// Setup HTTP routes
	http.SetupRoutes(app, itemHandler, userHandler, borrowingHandler)

	// Start the server
	err = app.Listen(fmt.Sprintf("%s:%s", config.Host, config.Port))
	if err != nil {
		log.Fatal("Error: ", err)
	}

	log.Println("Server is running on port", config.Port)
}
