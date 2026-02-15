package testutil

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// SetupTestDB creates an in-memory SQLite DB for tests.
// It enables foreign key constraints and verifies a connection can be established.
func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite3 DB: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping sqlite3 DB: %v", err)
	}
	// Ensure foreign keys are enforced to mimic real-world behavior.
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
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
