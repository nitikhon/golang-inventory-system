package seed

import (
	"log"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/util"
	"gorm.io/gorm"
)

// CheckIfDataExists checks if the database already contains data.
func CheckIfDataExists(db *gorm.DB) bool {
	var count int
	db.Raw("SELECT COUNT(*) FROM users").Scan(&count)
	return count > 0
}

func SeedDB(db *gorm.DB) error {
	if CheckIfDataExists(db) {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	log.Println("Seeding database...")

	crypto := util.NewAppCrptoUtil()

	if err := seedUsers(db, crypto); err != nil {
		log.Printf("An error occurs while trying to seed user data: %v", err)
	}
	if err := seedItems(db); err != nil {
		log.Printf("An error occurs while trying to seed user data: %v", err)
	}

	log.Println("Seeding database success")

	return nil
}

func seedUsers(db *gorm.DB, crpyto util.AppCryptoUtil) error {
	adminPassword, err := crpyto.HashPassword("admin123")
	if err != nil {
		return err
	}
	userPassword, err := crpyto.HashPassword("user123")
	if err != nil {
		return err
	}

	users := []entity.User{
		{
			Username:  "admin",
			Email:     "admin@example.com",
			Password:  adminPassword, // ✅ Hashed password
			Phone:     "0801234567",
			FirstName: "Admin",
			LastName:  "User",
			IsAdmin:   true,
		},
		{
			Username:  "john_doe",
			Email:     "john@example.com",
			Password:  userPassword, // ✅ Hashed password
			Phone:     "0987654321",
			FirstName: "John",
			LastName:  "Doe",
			IsAdmin:   false,
		},
		{
			Username:  "jane_smith",
			Email:     "jane@example.com",
			Password:  userPassword, // ✅ Hashed password
			Phone:     "0876543210",
			FirstName: "Jane",
			LastName:  "Smith",
			IsAdmin:   false,
		},
	}

	if err := db.Create(&users).Error; err != nil {
		return err
	}

	log.Println("Users seeded successfully")
	return nil
}

func seedItems(db *gorm.DB) error {
	items := []entity.Item{
		{
			Name:            "MacBook Pro 14",
			Description:     "Apple MacBook Pro 14-inch with M2 Pro chip",
			AvailableAmount: 3,
			TotalAmount:     5,
			Status:          "available",
		},
		{
			Name:            "iPad Air",
			Description:     "iPad Air with M1 chip, 64GB Wi-Fi",
			AvailableAmount: 8,
			TotalAmount:     10,
			Status:          "available",
		},
		{
			Name:            "Dell Monitor 27",
			Description:     "Dell 27-inch 4K USB-C Monitor",
			AvailableAmount: 2,
			TotalAmount:     4,
			Status:          "available",
		},
		{
			Name:            "Wireless Mouse",
			Description:     "Logitech MX Master 3S Wireless Mouse",
			AvailableAmount: 0,
			TotalAmount:     5,
			Status:          "borrowed",
		},
		{
			Name:            "USB-C Hub",
			Description:     "7-in-1 USB-C Hub with HDMI and Ethernet",
			AvailableAmount: 5,
			TotalAmount:     6,
			Status:          "available",
		},
	}

	// Batch create items
	if err := db.Create(&items).Error; err != nil {
		log.Println("Error seeding items:", err)
		return err
	}

	log.Println("Items seeded successfully!")
	return nil
}
