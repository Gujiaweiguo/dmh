package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMySQLTestConfig(t *testing.T) {
	config := GetMySQLTestConfig()
	assert.NotNil(t, config)
	assert.NotEmpty(t, config.Host)
	assert.NotEmpty(t, config.Port)
	assert.NotEmpty(t, config.User)
}

func TestGetEnv(t *testing.T) {
	// Test with default value
	value := getEnv("NON_EXISTENT_ENV_VAR", "default_value")
	assert.Equal(t, "default_value", value)
}

func TestGenerateTestDBName(t *testing.T) {
	name := generateTestDBName(t)
	assert.NotEmpty(t, name)
	assert.Contains(t, name, "t_")
}

func TestSkipIfNoMySQL(t *testing.T) {
	// This test verifies SkipIfNoMySQL doesn't panic
	// It will skip if MySQL is not available
	SkipIfNoMySQL(t)
	// If we get here, MySQL is available
	t.Log("MySQL is available for testing")
}

func TestSetupAndCleanupMySQLTestDB(t *testing.T) {
	SkipIfNoMySQL(t)

	// Test SetupMySQLTestDB
	db, dbName := SetupMySQLTestDB(t)
	assert.NotNil(t, db, "SetupMySQLTestDB should return non-nil DB")
	assert.NotEmpty(t, dbName, "SetupMySQLTestDB should return database name")
	assert.Contains(t, dbName, "t_", "Database name should have test prefix")

	// Verify DB is working
	var result int
	err := db.Raw("SELECT 1").Scan(&result).Error
	assert.NoError(t, err, "DB should be able to execute queries")
	assert.Equal(t, 1, result)

	// CleanupMySQLTestDB is called automatically via t.Cleanup
}

func TestGetEnv_Existing(t *testing.T) {
	// Test with existing env var (if set)
	t.Setenv("TEST_ENV_VAR", "test_value")
	value := getEnv("TEST_ENV_VAR", "default")
	assert.Equal(t, "test_value", value)
}
