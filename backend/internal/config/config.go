package config

import (
	"fmt"
	"os"

	"github.com/nitikhon/golang-inventory-system/internal/util"
)

// Config holds the configuration values for the application.
type Config struct {
	Host             string
	Port             string
	RedisHost        string
	RedisPort        int
	RateLimitBotMax  int
	RateLimitUserMax int
	DatabaseDSN      string
	Environment      string
}

// NewConfig creates a new Config instance and populates it with values from environment variables.
func NewConfig() Config {
	return Config{
		Host:             os.Getenv("HOST"),
		Port:             os.Getenv("PORT"),
		RedisHost:        os.Getenv("REDIS_HOST"),
		RedisPort:        util.ParseInt(os.Getenv("REDIS_PORT")),
		RateLimitBotMax:  util.ParseInt(os.Getenv("RATE_LIMIT_BOT_MAX")),
		RateLimitUserMax: util.ParseInt(os.Getenv("RATE_LIMIT_USER_MAX")),
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
