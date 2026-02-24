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
