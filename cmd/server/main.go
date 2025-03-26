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
	godotenv.Load()

	config := infrastructure.NewConfig()

	db, err := infrastructure.NewDatabase(*config)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	db.AutoMigrate(&domain.Item{})

	itemRepo := repository.NewItemRepository(db)
	itemService := app.NewItemService(itemRepo)
	itemHandler := http.NewItemHandler(itemService)

	app := fiber.New()

	http.SetupRoutes(app, itemHandler)

	err = app.Listen(fmt.Sprintf("%s:%s", config.Host, config.Port))
	if err != nil {
		log.Fatal("Error: ", err)
	}
	log.Println("Server is running on port", config.Port)
}
