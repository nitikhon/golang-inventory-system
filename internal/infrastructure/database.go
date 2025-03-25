package infrastructure

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDatabase initializes a new database connection using the provided configuration.
// It returns a pointer to the gorm.DB instance.
func NewDatabase(config Config) (*gorm.DB, error) {
	// Open a connection to the database using the provided DSN and custom logger.
	db, err := gorm.Open(postgres.Open(config.DatabaseDSN), &gorm.Config{
		Logger: newLogger(),
	})
	
	if err != nil {
		return nil, err
	}
	
	log.Println("Database connected")
	return db, nil
}
