package handlerutil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePathID(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		prefix   string
		errorMsg string
		expectID int64
		expectOK bool
	}{
		{
			name:     "valid ID",
			path:     "/api/v1/users/123",
			prefix:   "/api/v1/users/",
			errorMsg: "invalid user id",
			expectID: 123,
			expectOK: true,
		},
		{
			name:     "valid ID with trailing path",
			path:     "/api/v1/campaigns/456/details",
			prefix:   "/api/v1/campaigns/",
			errorMsg: "invalid campaign id",
			expectID: 456,
			expectOK: true,
		},
		{
			name:     "invalid ID - non-numeric",
			path:     "/api/v1/users/abc",
			prefix:   "/api/v1/users/",
			errorMsg: "invalid user id",
			expectID: 0,
			expectOK: false,
		},
		{
			name:     "ID is zero",
			path:     "/api/v1/users/0",
			prefix:   "/api/v1/users/",
			errorMsg: "invalid user id",
			expectID: 0,
			expectOK: true,
		},
		{
			name:     "negative ID",
			path:     "/api/v1/users/-1",
			prefix:   "/api/v1/users/",
			errorMsg: "invalid user id",
			expectID: -1,
			expectOK: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			id, ok := ParsePathID(req, w, tt.prefix, tt.errorMsg)

			assert.Equal(t, tt.expectID, id)
			assert.Equal(t, tt.expectOK, ok)
		})
	}
}

func TestParseQueryInt64(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		key          string
		defaultValue int64
		expectValue  int64
	}{
		{
			name:         "valid positive integer",
			query:        "?page=5",
			key:          "page",
			defaultValue: 1,
			expectValue:  5,
		},
		{
			name:         "missing key returns default",
			query:        "?other=value",
			key:          "page",
			defaultValue: 1,
			expectValue:  1,
		},
		{
			name:         "empty query returns default",
			query:        "",
			key:          "page",
			defaultValue: 1,
			expectValue:  1,
		},
		{
			name:         "non-numeric returns default",
			query:        "?page=abc",
			key:          "page",
			defaultValue: 1,
			expectValue:  1,
		},
		{
			name:         "zero returns default",
			query:        "?page=0",
			key:          "page",
			defaultValue: 1,
			expectValue:  1,
		},
		{
			name:         "negative returns default",
			query:        "?page=-5",
			key:          "page",
			defaultValue: 1,
			expectValue:  1,
		},
		{
			name:         "large positive integer",
			query:        "?id=9999999999",
			key:          "id",
			defaultValue: 0,
			expectValue:  9999999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test"+tt.query, nil)

			value := ParseQueryInt64(req, tt.key, tt.defaultValue)

			assert.Equal(t, tt.expectValue, value)
		})
	}
}

func TestParseQueryString(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		key         string
		expectValue string
	}{
		{
			name:        "valid string value",
			query:       "?name=test",
			key:         "name",
			expectValue: "test",
		},
		{
			name:        "empty string value",
			query:       "?name=",
			key:         "name",
			expectValue: "",
		},
		{
			name:        "missing key returns empty",
			query:       "?other=value",
			key:         "name",
			expectValue: "",
		},
		{
			name:        "no query returns empty",
			query:       "",
			key:         "name",
			expectValue: "",
		},
		{
			name:        "special characters",
			query:       "?name=test%20value",
			key:         "name",
			expectValue: "test value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test"+tt.query, nil)

			value := ParseQueryString(req, tt.key)

			assert.Equal(t, tt.expectValue, value)
		})
	}
}
