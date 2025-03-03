package agent

import (
	"os"
	"strconv"
)

// Settings holds all the configuration values
var (
	OrchestratorPort = getEnv("ORCHESTRATOR_PORT", "8080")
	ComputingPower   = getEnvInt("COMPUTING_POWER", 2)
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
