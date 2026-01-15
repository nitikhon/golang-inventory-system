package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/storage/redis/v3"
	"github.com/joho/godotenv"
	"github.com/nitikhon/golang-inventory-system/internal/adapter/inbound/http"
	"github.com/nitikhon/golang-inventory-system/internal/adapter/inbound/http/middleware"
	"github.com/nitikhon/golang-inventory-system/internal/adapter/outbound/database"
	"github.com/nitikhon/golang-inventory-system/internal/adapter/outbound/repository"
	"github.com/nitikhon/golang-inventory-system/internal/config"
	"github.com/nitikhon/golang-inventory-system/internal/core/service"
	"github.com/nitikhon/golang-inventory-system/internal/util"
)

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	// Initialize configuration
	config := config.NewConfig()

	// Connect to the database
	db, err := database.NewDatabase(config)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	if err := database.SetupDatabase(config, db); err != nil {
		log.Fatal("Error: ", err)
	}

	// Initialize repositories, services, and handlers
	itemRepo := repository.NewItemRepository(db)
	itemService := service.NewItemService(itemRepo)
	itemHandler := http.NewItemHandler(itemService)

	userRepo := repository.NewUserRepository(db)
	crypto := util.NewAppCrptoUtil()
	jwt := util.NewAppJWTUtil()
	userService := service.NewUserService(userRepo, crypto, jwt)
	userHandler := http.NewUserHandler(userService)

	borrowingRepo := repository.NewBorrowingRepository(db)
	borrowingService := service.NewBorrowingService(borrowingRepo, itemRepo, userRepo)
	borrowingHandler := http.NewBorrowingHandler(borrowingService)

	// Create a new Fiber app
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS, PATCH",
		AllowCredentials: true,
	}))

	// app.Static("/", "/app/dist")

	// Initialize Redis storage for rate limiting
	store := redis.New(redis.Config{
		Host: config.RedisHost,
		Port: config.RedisPort,
	})

	// Apply Rate Limiting Middleware Bot/DDOS for global
	app.Use(middleware.BotProtectionMiddleware(store, config.RateLimitBotMax))

	// Setup HTTP routes
	http.SetupRoutes(app, itemHandler, userHandler, borrowingHandler, middleware.RateLimitMiddleware(store, config.RateLimitUserMax))

	// app.Get("/*", func(c *fiber.Ctx) error {
	// 	return c.SendFile("/app/dist/index.html")
	// })

	go func() {
		for {
			time.Sleep(1 * time.Hour)
			borrowingService.MarkOverdueItems(context.Background())
		}
	}()

	// Start the server
	err = app.Listen(fmt.Sprintf("%s:%s", config.Host, config.Port))
	if err != nil {
		log.Fatal("Error: ", err)
	}

	log.Println("Server is running on port", config.Port)
}
