package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/nitikhon/golang-inventory-system/internal/infrastructure"
)

func main() {
	godotenv.Load()

	config := infrastructure.NewConfig()

	_, err := infrastructure.NewDatabase(*config)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	err = app.Listen(fmt.Sprintf("%s:%s", config.Host, config.Port))
	if err != nil {
		log.Fatal("Error: ", err)
	}
	log.Println("Server is running on port", config.Port)
}
