package main

import (
	"os"
	"testing"

	"dmh/api/internal/config"
)

func TestApplyEnvOverrides(t *testing.T) {
	os.Setenv("APP_HOST", "0.0.0.0")
	os.Setenv("APP_PORT", "9999")
	os.Setenv("JWT_SECRET", "jwt-test")
	os.Setenv("DB_HOST", "127.0.0.1")
	os.Setenv("DB_PORT", "3307")
	os.Setenv("DB_USER", "u1")
	os.Setenv("DB_PASSWORD", "p1")
	os.Setenv("DB_NAME", "d1")
	t.Cleanup(func() {
		os.Unsetenv("APP_HOST")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	})

	var c config.Config
	applyEnvOverrides(&c)

	if c.Host != "0.0.0.0" || c.Port != 9999 || c.Auth.AccessSecret != "jwt-test" {
		t.Fatalf("basic env overrides not applied")
	}
	if c.Mysql.DataSource == "" {
		t.Fatalf("db dsn should be generated")
	}
}
