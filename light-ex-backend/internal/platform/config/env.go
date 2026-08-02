package config

import (
	"os"
	"strconv"
	"time"
)

// getEnv returns the environment value if it exists,
// otherwise it returns the provided default value.
func getEnv(key, defaultValue string) string {

	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

// getEnvAsInt returns an integer environment variable
// or the provided default value.
func getEnvAsInt(key string, defaultValue int) int {

	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return number
}

// getEnvAsFloat returns a float64 environment variable
// or the provided default value.
func getEnvAsFloat(key string, defaultValue float64) float64 {

	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}

	return number
}

// getEnvAsDuration parses duration strings.
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {

	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return duration
}
