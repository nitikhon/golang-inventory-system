package database

import (
	"fmt"
	"log"
	"time"

	"github.com/nitikhon/golang-inventory-system/internal/config"
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDatabase initializes a new database connection using the provided configuration.
// It returns a pointer to the gorm.DB instance.
func NewDatabase(config config.Config) (*gorm.DB, error) {
	// Open a connection to the database using the provided DSN and custom logger.
	db, err := gorm.Open(postgres.Open(config.DatabaseDSN), &gorm.Config{
		Logger:                                   util.NewLogger(),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	log.Println("Database connected")
	return db, nil
}

func SetupDatabase(config config.Config, db *gorm.DB) error {
	switch config.Environment {
	case "development":
		log.Printf("%s: Dropping and recreating tables then seeding data\n", config.Environment)
		if err := db.Migrator().DropTable(&entity.Item{}, &entity.User{}, &entity.Borrowing{}); err != nil {
			return fmt.Errorf("failed to drop tables: %w", err)
		}

		if err := db.AutoMigrate(&entity.Item{}, &entity.User{}, &entity.Borrowing{}); err != nil {
			return fmt.Errorf("failed to auto-migrate: %w", err)
		}

		if err := util.SeedDB(db); err != nil {
			log.Println("Seeding warning: ", err)
		}

	case "production":
		log.Println("Production: Auto-migrating tables only")

		if err := db.AutoMigrate(&entity.Item{}, &entity.User{}, &entity.Borrowing{}); err != nil {
			return fmt.Errorf("failed to auto-migrate: %w", err)
		}
	default:
		log.Println("Unknown environment: Auto-migrating tables only")
	}
	return nil
}
