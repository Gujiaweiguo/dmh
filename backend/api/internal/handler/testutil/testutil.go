package testutil

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const testDSN = "root:Admin168@tcp(127.0.0.1:3306)/dmh_test?charset=utf8mb4&parseTime=true&loc=Local&multiStatements=true"

// SetupTestDB creates a *sql.DB connection for testing (legacy, uses database/sql).
func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("mysql", testDSN)
	if err != nil {
		t.Fatalf("failed to open MySQL test DB: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping MySQL DB: %v", err)
	}
	return db
}

// SetupGormTestDB creates a *gorm.DB connection for testing with relaxed sql_mode
// to avoid "Invalid default value for 'created_at'" errors on MySQL 8.0+.
func SetupGormTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(mysql.Open(testDSN), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open MySQL test DB with GORM: %v", err)
	}

	// Set relaxed sql_mode to allow CURRENT_TIMESTAMP default on datetime columns
	// This removes NO_ZERO_IN_DATE, NO_ZERO_DATE which cause issues with GORM's
	// DEFAULT CURRENT_TIMESTAMP syntax on MySQL 8.0+
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get underlying sql.DB: %v", err)
	}
	if _, err := sqlDB.Exec("SET SESSION sql_mode = 'NO_ENGINE_SUBSTITUTION'"); err != nil {
		t.Fatalf("failed to set sql_mode: %v", err)
	}

	return db
}

// MakeRequest creates an HTTP request for testing handlers.
// If a non-nil body is provided, Content-Type is assumed to be application/json.
func MakeRequest(method, path string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		// In tests, returning nil would complicate; panic to surface misconfig quickly.
		panic(err)
	}
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

// ExecuteRequest executes the given handler with the provided request and returns the response recorder.
func ExecuteRequest(handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// GenUniquePhone generates a unique phone number for testing.
func GenUniquePhone() string {
	return fmt.Sprintf("138%08d", time.Now().UnixNano()%100000000)
}

// GenUniqueUsername generates a unique username for testing.
func GenUniqueUsername(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

// GenUniqueUnionID generates a unique WeChat unionid for testing.
func GenUniqueUnionID() string {
	return fmt.Sprintf("union_%d", time.Now().UnixNano())
}

// GenUniqueCode generates a unique code for testing (e.g., menu code, role code).
func GenUniqueCode(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

// WithTransaction runs a test function within a database transaction that is rolled back afterward.
// This ensures test data isolation without needing separate databases.
func WithTransaction(t *testing.T, db *gorm.DB, fn func(tx *gorm.DB)) {
	t.Helper()
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	fn(tx)

	if err := tx.Rollback().Error; err != nil && err != gorm.ErrInvalidTransaction {
		t.Logf("warning: rollback error: %v", err)
	}
}
