package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/nitikhon/golang-inventory-system/internal/app"
	"github.com/nitikhon/golang-inventory-system/internal/domain"
	"github.com/nitikhon/golang-inventory-system/internal/infrastructure"
	"github.com/nitikhon/golang-inventory-system/internal/infrastructure/repository"
	"github.com/nitikhon/golang-inventory-system/internal/interfaces/http"
)

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	// Initialize configuration
	config := infrastructure.NewConfig()

	// Connect to the database
	db, err := infrastructure.NewDatabase(*config)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	// Auto-migrate database models
	db.AutoMigrate(&domain.Item{}, &domain.User{})

	// Seed the database if it's empty
	if err := infrastructure.SeedDB(db); err != nil {
		log.Println("Error: ", err)
	} else {
		log.Println("Database is empty, seeding...")
	}

	// Initialize repositories, services, and handlers
	itemRepo := repository.NewItemRepository(db)
	itemService := app.NewItemService(itemRepo)
	itemHandler := http.NewItemHandler(itemService)

	userRepo := repository.NewUserRepository(db)
	userService := app.NewUserService(userRepo)
	userHandler := http.NewUserHandler(userService)

	// Create a new Fiber app
	app := fiber.New()

	// Setup HTTP routes
	http.SetupRoutes(app, itemHandler, userHandler)

	// Start the server
	err = app.Listen(fmt.Sprintf("%s:%s", config.Host, config.Port))
	if err != nil {
		log.Fatal("Error: ", err)
	}

	log.Println("Server is running on port", config.Port)
}
