package util

import (
	"log"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
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

	crypto := NewAppCrptoUtil()

	if err := seedUsers(db, crypto); err != nil {
		log.Printf("An error occurs while trying to seed user data: %v", err)
	}
	if err := seedItems(db); err != nil {
		log.Printf("An error occurs while trying to seed user data: %v", err)
	}
	if err := seedBorrowings(db); err != nil {
		log.Printf("An error occurs while trying to seed borrowing data: %v", err)
	}

	log.Println("Seeding database success")

	return nil
}

func seedUsers(db *gorm.DB, crpyto AppCryptoUtil) error {
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
		{
			Name:            "Projector 4K",
			Description:     "Sony 4K Home Theater Projector",
			AvailableAmount: 0,
			TotalAmount:     1,
			Status:          "maintenance",
		},
		{
			Name:            "Old Laptop",
			Description:     "Lenovo ThinkPad X220",
			AvailableAmount: 0,
			TotalAmount:     1,
			Status:          "lost",
		},
		{
			Name:            "iPhone 15 Pro",
			Description:     "Apple iPhone 15 Pro 256GB - Blue Titanium",
			AvailableAmount: 2,
			TotalAmount:     2,
			Status:          "available",
		},
		{
			Name:            "Mechanical Keyboard",
			Description:     "Keychron K2 Pro Wireless Mechanical Keyboard",
			AvailableAmount: 4,
			TotalAmount:     5,
			Status:          "available",
		},
		{
			Name:            "GoPro Hero 11",
			Description:     "GoPro Hero 11 Black Action Camera",
			AvailableAmount: 1,
			TotalAmount:     3,
			Status:          "available",
		},
		{
			Name:            "Sony WH-1000XM5",
			Description:     "Wireless Noise Canceling Headphones",
			AvailableAmount: 5,
			TotalAmount:     5,
			Status:          "available",
		},
		{
			Name:            "Nintendo Switch OLED",
			Description:     "Nintendo Switch Console with White Joy-Con",
			AvailableAmount: 1,
			TotalAmount:     3,
			Status:          "available",
		},
		{
			Name:            "Arduino Starter Kit",
			Description:     "Official Arduino Starter Kit for beginners",
			AvailableAmount: 10,
			TotalAmount:     10,
			Status:          "available",
		},
		{
			Name:            "Raspberry Pi 4",
			Description:     "Raspberry Pi 4 Model B 8GB RAM",
			AvailableAmount: 8,
			TotalAmount:     15,
			Status:          "available",
		},
		{
			Name:            "Standing Desk Converter",
			Description:     "Adjustable height desk riser",
			AvailableAmount: 0,
			TotalAmount:     2,
			Status:          "borrowed",
		},
		{
			Name:            "Ergonomic Chair",
			Description:     "Mesh office chair with lumbar support",
			AvailableAmount: 2,
			TotalAmount:     10,
			Status:          "available",
		},
		{
			Name:            "Blue Yeti Microphone",
			Description:     "USB Microphone for Recording and Streaming",
			AvailableAmount: 3,
			TotalAmount:     4,
			Status:          "available",
		},
		{
			Name:            "Logitech C920 Webcam",
			Description:     "HD Pro Webcam for video calling",
			AvailableAmount: 6,
			TotalAmount:     8,
			Status:          "available",
		},
		{
			Name:            "DJI Mini 3 Drone",
			Description:     "Lightweight Camera Drone with 4K Video",
			AvailableAmount: 0,
			TotalAmount:     1,
			Status:          "maintenance",
		},
		{
			Name:            "Meta Quest 3",
			Description:     "Advanced All-in-One VR Headset",
			AvailableAmount: 2,
			TotalAmount:     2,
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

func seedBorrowings(db *gorm.DB) error {
	borrowings := []entity.Borrowing{
		// Pending
		{
			UserID:          2, // john_doe
			ItemID:          1, // MacBook Pro 14
			Description:     "Need for project presentation",
			BorrowedAt:      "",
			DueDate:         "2023-12-31T00:00:00Z",
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_PENDING,
			ApprovalStatus:  entity.APPROVAL_PENDING,
		},
		// Active
		{
			UserID:          2, // john_doe
			ItemID:          2, // iPad Air
			Description:     "Design work",
			BorrowedAt:      "2023-10-01T10:00:00Z",
			DueDate:         "2023-10-15T18:00:00Z",
			ApprovedAt:      "2023-10-01T09:00:00Z",
			ApprovedBy:      1,
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_ACTIVE,
			ApprovalStatus:  entity.APPROVAL_APPROVED,
		},
		// Returned
		{
			UserID:          2, // john_doe
			ItemID:          3, // Dell Monitor
			Description:     "Extra screen",
			BorrowedAt:      "2023-09-01T09:00:00Z",
			ReturnedAt:      "2023-09-05T17:00:00Z",
			DueDate:         "2023-09-10T18:00:00Z",
			ApprovedAt:      "2023-09-01T08:30:00Z",
			ApprovedBy:      1,
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_RETURNED,
			ApprovalStatus:  entity.APPROVAL_APPROVED,
		},
		// Overdue
		{
			UserID:          2, // john_doe
			ItemID:          5, // USB-C Hub
			Description:     "Connectivity",
			BorrowedAt:      "2023-08-01T09:00:00Z",
			DueDate:         "2023-08-05T18:00:00Z",
			ApprovedAt:      "2023-08-01T08:30:00Z",
			ApprovedBy:      1,
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_OVERDUE,
			ApprovalStatus:  entity.APPROVAL_APPROVED,
		},
		// Rejected
		{
			UserID:          2, // john_doe
			ItemID:          4, // Wireless Mouse
			Description:     "Gaming",
			DueDate:         "2023-11-01T00:00:00Z",
			RejectedAt:      "2023-10-25T10:00:00Z",
			RejectedBy:      1,
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_CANCELLED,
			ApprovalStatus:  entity.APPROVAL_REJECTED,
		},
		// Lost
		{
			UserID:          2, // john_doe
			ItemID:          2, // iPad Air
			Description:     "Lost during commute",
			BorrowedAt:      "2023-07-01T09:00:00Z",
			DueDate:         "2023-07-10T18:00:00Z",
			ApprovedAt:      "2023-07-01T08:30:00Z",
			ApprovedBy:      1,
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_LOST,
			ApprovalStatus:  entity.APPROVAL_APPROVED,
		},
		// More mock data for User 2 (John)
		{
			UserID:          2, // john_doe
			ItemID:          8, // iPhone 15 Pro
			Description:     "Mobile testing",
			BorrowedAt:      "2023-11-05T10:00:00Z",
			DueDate:         "2023-11-20T18:00:00Z",
			ApprovedAt:      "2023-11-05T09:30:00Z",
			ApprovedBy:      1,
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_ACTIVE,
			ApprovalStatus:  entity.APPROVAL_APPROVED,
		},
		{
			UserID:          2, // john_doe
			ItemID:          9, // Mechanical Keyboard
			Description:     "Coding session",
			BorrowedAt:      "",
			DueDate:         "2023-12-05T18:00:00Z",
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_PENDING,
			ApprovalStatus:  entity.APPROVAL_PENDING,
		},
		{
			UserID:          2,  // john_doe
			ItemID:          13, // Arduino Starter Kit
			Description:     "IoT Learning",
			BorrowedAt:      "2023-06-01T09:00:00Z",
			ReturnedAt:      "2023-06-20T15:00:00Z",
			DueDate:         "2023-06-30T18:00:00Z",
			ApprovedAt:      "2023-06-01T08:30:00Z",
			ApprovedBy:      1,
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_RETURNED,
			ApprovalStatus:  entity.APPROVAL_APPROVED,
		},

		// Mock data for User 3 (Jane)
		{
			UserID:          3, // jane_smith
			ItemID:          1, // MacBook Pro 14
			Description:     "Main work laptop",
			BorrowedAt:      "2023-11-01T09:00:00Z",
			DueDate:         "2023-11-30T18:00:00Z",
			ApprovedAt:      "2023-11-01T08:30:00Z",
			ApprovedBy:      1,
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_ACTIVE,
			ApprovalStatus:  entity.APPROVAL_APPROVED,
		},
		{
			UserID:          3,  // jane_smith
			ItemID:          11, // Sony WH-1000XM5
			Description:     "Focus time",
			BorrowedAt:      "2023-11-10T09:00:00Z",
			DueDate:         "2023-11-17T18:00:00Z",
			ApprovedAt:      "2023-11-10T08:45:00Z",
			ApprovedBy:      1,
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_ACTIVE,
			ApprovalStatus:  entity.APPROVAL_APPROVED,
		},
		{
			UserID:          3,  // jane_smith
			ItemID:          14, // Raspberry Pi 4
			Description:     "Home server setup",
			BorrowedAt:      "",
			DueDate:         "2023-12-15T18:00:00Z",
			BorrowingAmount: 2,
			BorrowingStatus: entity.BORROWING_PENDING,
			ApprovalStatus:  entity.APPROVAL_PENDING,
		},
		{
			UserID:          3,  // jane_smith
			ItemID:          12, // Nintendo Switch OLED
			Description:     "Team building activity",
			BorrowedAt:      "2023-10-01T10:00:00Z",
			DueDate:         "2023-10-05T18:00:00Z",
			ApprovedAt:      "2023-10-01T09:30:00Z",
			ApprovedBy:      1,
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_OVERDUE,
			ApprovalStatus:  entity.APPROVAL_APPROVED,
		},
		{
			UserID:          3,  // jane_smith
			ItemID:          10, // GoPro Hero 11
			Description:     "Vlog project",
			BorrowedAt:      "2023-09-15T09:00:00Z",
			ReturnedAt:      "2023-09-20T14:00:00Z",
			DueDate:         "2023-09-22T18:00:00Z",
			ApprovedAt:      "2023-09-15T08:30:00Z",
			ApprovedBy:      1,
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_RETURNED,
			ApprovalStatus:  entity.APPROVAL_APPROVED,
		},
		{
			UserID:          3,  // jane_smith
			ItemID:          16, // Ergonomic Chair
			Description:     "Back pain relief",
			BorrowedAt:      "2023-08-01T09:00:00Z",
			DueDate:         "2023-12-31T18:00:00Z",
			ApprovedAt:      "2023-08-01T08:30:00Z",
			ApprovedBy:      1,
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_ACTIVE,
			ApprovalStatus:  entity.APPROVAL_APPROVED,
		},
		{
			UserID:          3,  // jane_smith
			ItemID:          17, // Blue Yeti Microphone
			Description:     "Podcast recording",
			BorrowedAt:      "",
			DueDate:         "2023-11-01T18:00:00Z",
			RejectedAt:      "2023-10-31T10:00:00Z",
			RejectedBy:      1,
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_CANCELLED,
			ApprovalStatus:  entity.APPROVAL_REJECTED,
		},
		{
			UserID:          3, // jane_smith
			ItemID:          6, // Projector 4K
			Description:     "Movie night",
			BorrowedAt:      "2023-07-20T17:00:00Z",
			ReturnedAt:      "2023-07-21T09:00:00Z",
			DueDate:         "2023-07-22T12:00:00Z",
			ApprovedAt:      "2023-07-20T16:00:00Z",
			ApprovedBy:      1,
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_RETURNED,
			ApprovalStatus:  entity.APPROVAL_APPROVED,
		},
		{
			UserID:          3,  // jane_smith
			ItemID:          18, // Logitech Webcam
			Description:     "Zoom meetings",
			BorrowedAt:      "2023-06-15T09:00:00Z",
			ReturnedAt:      "2023-06-18T17:00:00Z",
			DueDate:         "2023-06-20T18:00:00Z",
			ApprovedAt:      "2023-06-15T08:30:00Z",
			ApprovedBy:      1,
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_RETURNED,
			ApprovalStatus:  entity.APPROVAL_APPROVED,
		},
	}

	if err := db.Create(&borrowings).Error; err != nil {
		log.Println("Error seeding borrowings:", err)
		return err
	}

	log.Println("Borrowings seeded successfully!")
	return nil
}
