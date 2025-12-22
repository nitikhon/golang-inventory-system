package util

import (
	"os"

	"github.com/joho/godotenv"
)

// If the environment variable is not set, it panics with an error message.
func GetEnvOrPanic(key string) []byte {
	godotenv.Load()
	val := os.Getenv(key)
	if val == "" {
		panic("missing env var: " + key)
	}
	return []byte(val)
}