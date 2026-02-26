// Integration test environment configuration
package testutil

import (
	"os"
	"strings"
)

// IntegrationConfig holds configuration for integration tests
type IntegrationConfig struct {
	BaseURL       string
	AdminUsername string
	AdminPassword string
}

// GetIntegrationConfig returns integration test configuration from environment
func GetIntegrationConfig() *IntegrationConfig {
	return &IntegrationConfig{
		BaseURL:       getEnv("DMH_INTEGRATION_BASE_URL", "http://localhost:8889"),
		AdminUsername: getEnv("DMH_TEST_ADMIN_USERNAME", "admin"),
		AdminPassword: getEnv("DMH_TEST_ADMIN_PASSWORD", "123456"),
	}
}

// SkipReason represents a reason for skipping a test
type SkipReason string

const (
	// SkipReasonAPIUnavailable API service is not available
	SkipReasonAPIUnavailable SkipReason = "API_UNAVAILABLE"
	// SkipReasonMySQLUnavailable MySQL database is not available
	SkipReasonMySQLUnavailable SkipReason = "MYSQL_UNAVAILABLE"
	// SkipReasonRedisUnavailable Redis is not available
	SkipReasonRedisUnavailable SkipReason = "REDIS_UNAVAILABLE"
	// SkipReasonLoginFailed Login to API failed
	SkipReasonLoginFailed SkipReason = "LOGIN_FAILED"
	// SkipReasonDataPrepFailed Test data preparation failed
	SkipReasonDataPrepFailed SkipReason = "DATA_PREP_FAILED"
)

// getEnv returns environment variable value or default
func getEnv(key, defaultValue string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return defaultValue
}
