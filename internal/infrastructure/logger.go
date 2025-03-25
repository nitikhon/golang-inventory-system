package infrastructure

import (
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"
)

// newLogger creates a new logger for detailed SQL logging
func newLogger() logger.Interface {
	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	return newLogger
}
