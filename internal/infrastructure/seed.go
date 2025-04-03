package infrastructure

import (
	"errors"
	"log"
	"os"

	"gorm.io/gorm"
)

// CheckIfDataExists checks if the database already contains data.
func CheckIfDataExists(db *gorm.DB) bool {
	var count int
	db.Raw("SELECT COUNT(*) FROM users").Scan(&count)
	return count > 0
}

// SeedDB seeds the database with initial data if it is empty.
func SeedDB(db *gorm.DB) error {
	if CheckIfDataExists(db) {
		return errors.New("database already seeded")
	}

	log.Println("Seeding database...")

	// Read the SQL file
	sqlFile := "internal/infrastructure/seed.sql"
	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Println("Error reading seed file:", err)
		return err
	}

	// Execute the SQL commands
	if err := db.Exec(string(sqlBytes)).Error; err != nil {
		log.Println("Error executing seed file:", err)
		return err
	} else {
		log.Println("Database seeded successfully!")
	}

	return nil
}
