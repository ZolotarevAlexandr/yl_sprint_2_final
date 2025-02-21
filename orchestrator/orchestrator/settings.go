package orchestrator

import (
	"os"
	"strconv"
)

// Settings holds all the configuration values
var (
	Port                 = getEnv("ORCHESTRATOR_PORT", "8080")
	AdditionTimeMs       = getEnvInt("TIME_ADDITION_MS", 1000)
	SubtractionTimeMs    = getEnvInt("TIME_SUBTRACTION_MS", 1000)
	MultiplicationTimeMs = getEnvInt("TIME_MULTIPLICATIONS_MS", 1000)
	DivisionTimeMs       = getEnvInt("TIME_DIVISIONS_MS", 1000)
)

// getEnv retrieves a string environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt retrieves an integer environment variable or returns a default value.
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
