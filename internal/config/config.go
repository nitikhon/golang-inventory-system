package config

import (
	"fmt"
	"os"
)

// Config holds the configuration values for the application.
type Config struct {
	Host        string
	Port        string
	DatabaseDSN string
	Environment string
}

// NewConfig creates a new Config instance and populates it with values from environment variables.
func NewConfig() Config {
	return Config{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
		DatabaseDSN: fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT")),
		Environment: os.Getenv("ENVIRONMENT"),
	}
}
