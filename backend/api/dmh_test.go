package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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

func TestApplyEnvOverrides_PartialDBVars(t *testing.T) {
	os.Setenv("DB_HOST", "192.168.1.1")
	os.Setenv("DB_PORT", "3308")
	t.Cleanup(func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	})

	var c config.Config
	applyEnvOverrides(&c)

	if c.Mysql.DataSource == "" {
		t.Fatalf("db dsn should be generated with partial env vars")
	}
}

func TestApplyEnvOverrides_EmptyEnv(t *testing.T) {
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

	os.Unsetenv("APP_HOST")
	os.Unsetenv("APP_PORT")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")

	var c config.Config
	c.Host = "original-host"
	c.Port = 8888
	originalDSN := "user:pass@tcp(localhost:3306)/db"
	c.Mysql.DataSource = originalDSN

	applyEnvOverrides(&c)

	if c.Host != "original-host" {
		t.Fatalf("host should not change when env vars are empty")
	}
	if c.Mysql.DataSource != originalDSN {
		t.Fatalf("datasource should not change when DB env vars are empty")
	}
}

func TestApplyEnvOverrides_InvalidPort(t *testing.T) {
	os.Setenv("APP_PORT", "invalid")
	os.Setenv("DB_HOST", "127.0.0.1")
	t.Cleanup(func() {
		os.Unsetenv("APP_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	})

	var c config.Config
	c.Port = 8888
	applyEnvOverrides(&c)

	if c.Port != 8888 {
		t.Fatalf("port should not change when APP_PORT is invalid")
	}
}

func TestApplyEnvOverrides_NegativePort(t *testing.T) {
	os.Setenv("APP_PORT", "-1")
	os.Setenv("DB_HOST", "127.0.0.1")
	t.Cleanup(func() {
		os.Unsetenv("APP_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	})

	var c config.Config
	c.Port = 8888
	applyEnvOverrides(&c)

	if c.Port != 8888 {
		t.Fatalf("port should not change when APP_PORT is negative")
	}
}

func TestApplyEnvOverrides_ZeroPort(t *testing.T) {
	os.Setenv("APP_PORT", "0")
	os.Setenv("DB_HOST", "127.0.0.1")
	t.Cleanup(func() {
		os.Unsetenv("APP_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	})

	var c config.Config
	c.Port = 8888
	applyEnvOverrides(&c)

	if c.Port != 8888 {
		t.Fatalf("port should not change when APP_PORT is zero")
	}
}

func TestServePosterFile_EmptyFilename(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/posters/", nil)
	rec := httptest.NewRecorder()

	servePosterFile(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestServePosterFile_FileNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/posters/nonexistent.png", nil)
	rec := httptest.NewRecorder()

	servePosterFile(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 for non-existent file, got %d", rec.Code)
	}
}

func TestServePosterFile_ExtractsFilename(t *testing.T) {
	testCases := []struct {
		path     string
		expected string
	}{
		{"/api/v1/posters/test.png", "test.png"},
		{"/api/v1/posters/image.jpg", "image.jpg"},
		{"/posters/file.png", "file.png"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			pathParts := strings.Split(tc.path, "/")
			filename := pathParts[len(pathParts)-1]
			if filename != tc.expected {
				t.Errorf("path %s: expected %s, got %s", tc.path, tc.expected, filename)
			}
		})
	}
}

func TestApplyEnvOverrides_OnlyPassword(t *testing.T) {
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_NAME")
	os.Setenv("DB_PASSWORD", "secret123")
	t.Cleanup(func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	})

	var c config.Config
	applyEnvOverrides(&c)

	if c.Mysql.DataSource == "" {
		t.Fatalf("db dsn should be generated with only password")
	}
}

func TestApplyEnvOverrides_JWTSecretWhitespace(t *testing.T) {
	os.Setenv("JWT_SECRET", "  trimmed-secret  ")
	os.Setenv("DB_HOST", "127.0.0.1")
	t.Cleanup(func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	})

	var c config.Config
	applyEnvOverrides(&c)

	// JWT secret IS trimmed (using strings.TrimSpace)
	if c.Auth.AccessSecret != "trimmed-secret" {
		t.Fatalf("JWT secret should be trimmed, got: %s", c.Auth.AccessSecret)
	}
}

func TestServePosterFile_Success(t *testing.T) {
	// Create test directory and file
	testDir := "/opt/data/posters"
	testFile := filepath.Join(testDir, "test-success.png")

	// Ensure directory exists
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create a minimal test file (doesn't need to be a valid PNG for this test)
	testContent := []byte("fake-png-content-for-testing")
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posters/test-success.png", nil)
	rec := httptest.NewRecorder()

	servePosterFile(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", rec.Header().Get("Content-Type"))
	}
}
