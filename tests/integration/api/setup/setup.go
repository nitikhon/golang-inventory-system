package setup

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/nitikhon/golang-inventory-system/internal/adapter/inbound/http"
	"github.com/nitikhon/golang-inventory-system/internal/adapter/outbound/repository"
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/core/service"
	"github.com/nitikhon/golang-inventory-system/internal/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestServer struct {
	App    *fiber.App
	DB     *gorm.DB
	Config *TestConfig
	t      *testing.T
}

type TestConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// Enable verbose logging for tests
var verboseLogging = os.Getenv("TEST_VERBOSE") == "true"

func init() {
	if err := loadEnvFile(); err != nil {
		log.Printf("[TEST SETUP] Warning: Could not load .env file: %v\n", err)
	}
}

func loadEnvFile() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	for range 5 {
		envPath := filepath.Join(cwd, ".env")
		if _, err := os.Stat(envPath); err == nil {
			log.Printf("[TEST SETUP] Loading .env from: %s\n", envPath)
			return godotenv.Load(envPath)
		}
		cwd = filepath.Dir(cwd)
	}

	return fmt.Errorf(".env file not found")
}

func getEnvOrFail(t *testing.T, key string) string {
	value := os.Getenv(key)
	if value == "" {
		t.Fatalf("[TEST SETUP] Missing required environment variable: %s", key)
	}
	return value
}

func NewTestServer(t *testing.T) *TestServer {
	t.Log("[TEST SETUP] Initializing test server...")

	config := &TestConfig{
		DBHost:     getEnvOrFail(t, "TEST_DB_HOST"),
		DBPort:     getEnvOrFail(t, "TEST_DB_PORT"),
		DBUser:     getEnvOrFail(t, "TEST_DB_USER"),
		DBPassword: getEnvOrFail(t, "TEST_DB_PASSWORD"),
		DBName:     getEnvOrFail(t, "TEST_DB_NAME"),
	}

	t.Logf("[TEST SETUP] Connecting to database: %s:%s/%s", config.DBHost, config.DBPort, config.DBName)

	db := setupTestDatabase(t, config)

	t.Log("[TEST SETUP] Setting up repositories...")
	// Setup repositories
	itemRepo := repository.NewItemRepository(db)
	userRepo := repository.NewUserRepository(db)
	borrowingRepo := repository.NewBorrowingRepository(db)

	t.Log("[TEST SETUP] Setting up services...")
	// Setup services
	crypto := util.NewAppCrptoUtil()
	jwt := util.NewAppJWTUtil()
	itemService := service.NewItemService(itemRepo)
	userService := service.NewUserService(userRepo, crypto, jwt)
	borrowingService := service.NewBorrowingService(borrowingRepo, itemRepo, userRepo)

	t.Log("[TEST SETUP] Setting up handlers...")
	// Setup handlers
	itemHandler := http.NewItemHandler(itemService)
	userHandler := http.NewUserHandler(userService)
	borrowingHandler := http.NewBorrowingHandler(borrowingService)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Printf("[TEST SERVER] Error: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Setup routes
	http.SetupRoutes(app, itemHandler, userHandler, borrowingHandler)

	t.Log("[TEST SETUP] Test server initialized successfully!")

	return &TestServer{
		App:    app,
		DB:     db,
		Config: config,
		t:      t,
	}
}

func (s *TestServer) Cleanup() {
	s.t.Log("[TEST CLEANUP] Cleaning up test data...")

	// Clean up test data
	if err := s.DB.Exec("TRUNCATE TABLE borrowings CASCADE").Error; err != nil {
		s.t.Logf("[TEST CLEANUP] Warning: Failed to truncate borrowings: %v", err)
	}
	if err := s.DB.Exec("TRUNCATE TABLE items CASCADE").Error; err != nil {
		s.t.Logf("[TEST CLEANUP] Warning: Failed to truncate items: %v", err)
	}
	if err := s.DB.Exec("TRUNCATE TABLE users CASCADE").Error; err != nil {
		s.t.Logf("[TEST CLEANUP] Warning: Failed to truncate users: %v", err)
	}

	// Close database connection
	sqlDB, err := s.DB.DB()
	if err != nil {
		s.t.Logf("[TEST CLEANUP] Warning: Failed to get SQL DB: %v", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		s.t.Logf("[TEST CLEANUP] Warning: Failed to close database connection: %v", err)
	}

	s.t.Log("[TEST CLEANUP] Cleanup completed!")
}

func setupTestDatabase(t *testing.T, config *TestConfig) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)

	// Configure GORM logger based on verbosity
	var gormLogger logger.Interface
	if verboseLogging {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	t.Log("[TEST SETUP] Opening database connection...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   gormLogger,
	})
	if err != nil {
		t.Fatalf("[TEST SETUP] Failed to connect to test database: %v", err)
	}

	t.Log("[TEST SETUP] Running migrations...")
	err = db.AutoMigrate(&entity.Item{}, &entity.User{}, &entity.Borrowing{})
	if err != nil {
		t.Fatalf("[TEST SETUP] Failed to migrate test tables: %v", err)
	}

	t.Log("[TEST SETUP] Cleaning existing data...")
	db.Exec("TRUNCATE TABLE borrowings CASCADE")
	db.Exec("TRUNCATE TABLE items RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")

	t.Log("[TEST SETUP] Seeding test data...")
	if err := seedTestData(t, db); err != nil {
		t.Fatalf("[TEST SETUP] Failed to seed test data: %v", err)
	}

	t.Log("[TEST SETUP] Database setup completed!")
	return db
}

func seedTestData(t *testing.T, db *gorm.DB) error {
	crypto := util.NewAppCrptoUtil()
	hashedPw, err := crypto.HashPassword("P@ssw0rd")
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	t.Log("[TEST SETUP] Creating test users...")
	adminUser := &entity.User{
		Username:  "test_admin",
		Email:     "admin@test.com",
		Password:  hashedPw,
		Phone:     "1234567890",
		FirstName: "Test",
		LastName:  "Admin",
		IsAdmin:   true,
	}

	regularUser := &entity.User{
		Username:  "test_user",
		Email:     "user@test.com",
		Password:  hashedPw,
		Phone:     "0987654321",
		FirstName: "Test",
		LastName:  "User",
		IsAdmin:   false,
	}

	if err := db.Create([]*entity.User{adminUser, regularUser}).Error; err != nil {
		return fmt.Errorf("failed to create test users: %w", err)
	}
	t.Logf("[TEST SETUP] Created users: %s (admin), %s (regular)", adminUser.Username, regularUser.Username)

	t.Log("[TEST SETUP] Creating test items...")
	items := []entity.Item{
		{
			Name:            "office chair ergonomic",
			Description:     "Ergonomic office chair with lumbar support and adjustable height",
			AvailableAmount: 4,
			TotalAmount:     6,
			Status:          "available",
		},
		{
			Name:            "standing desk electric",
			Description:     "Height-adjustable electric standing desk 120x60cm",
			AvailableAmount: 2,
			TotalAmount:     3,
			Status:          "available",
		},
		{
			Name:            "wireless keyboard",
			Description:     "Mechanical wireless keyboard with RGB backlight",
			AvailableAmount: 7,
			TotalAmount:     10,
			Status:          "available",
		},
		{
			Name:            "conference camera 4k",
			Description:     "4K webcam with auto-focus for video conferences",
			AvailableAmount: 0,
			TotalAmount:     4,
			Status:          "borrowed",
		},
		{
			Name:            "bluetooth headset",
			Description:     "Noise-cancelling wireless headset for calls and meetings",
			AvailableAmount: 6,
			TotalAmount:     8,
			Status:          "available",
		},
		{
			Name:            "portable projector",
			Description:     "Mini LED projector 1080p with wireless connectivity",
			AvailableAmount: 1,
			TotalAmount:     2,
			Status:          "maintenance",
		},
		{
			Name:            "external hard drive",
			Description:     "1TB USB 3.0 external hard drive for data backup",
			AvailableAmount: 3,
			TotalAmount:     5,
			Status:          "available",
		},
	}

	if err := db.Create(&items).Error; err != nil {
		return fmt.Errorf("failed to create test items: %w", err)
	}
	t.Logf("[TEST SETUP] Created %d test items", len(items))

	t.Log("[TEST SETUP] Creating test borrowings...")

	users := []entity.User{*adminUser, *regularUser}

	borrowings := []entity.Borrowing{
		{
			UserID:          findUserID(users, "test_user"),
			ItemID:          findItemID(items, "office chair ergonomic"),
			Description:     "Need for home office",
			BorrowedAt:      "2023-01-01T10:00:00Z",
			DueDate:         "2023-02-01T10:00:00Z",
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_ACTIVE,
			ApprovalStatus:  entity.APPROVAL_APPROVED,
			ApprovedAt:      "2023-01-01T10:05:00Z",
			ApprovedBy:      findUserID(users, "test_admin"),
		},
		{
			UserID:          findUserID(users, "test_user"),
			ItemID:          findItemID(items, "wireless keyboard"),
			Description:     "Old one broke",
			BorrowedAt:      "2023-01-15T14:30:00Z",
			DueDate:         "2023-02-15T14:30:00Z",
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_PENDING,
			ApprovalStatus:  entity.APPROVAL_PENDING,
		},
	}

	if err := db.Create(&borrowings).Error; err != nil {
		return fmt.Errorf("failed to create test borrowings: %w", err)
	}
	t.Logf("[TEST SETUP] Created %d test borrowings", len(borrowings))

	return nil
}

func findUserID(users []entity.User, username string) uint {
	for _, u := range users {
		if u.Username == username {
			return u.ID
		}
	}
	return 0
}

func findItemID(items []entity.Item, name string) uint {
	for _, i := range items {
		if i.Name == name {
			return i.ID
		}
	}
	return 0
}
