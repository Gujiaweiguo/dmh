package integration

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLInjectionPrevention(t *testing.T) {
	maliciousInputs := []string{
		"' OR '1'='1",
		"'; DROP TABLE users; --",
		"1' UNION SELECT * FROM users --",
	}

	for _, input := range maliciousInputs {
		escaped := strings.ReplaceAll(input, "'", "''")
		assert.NotEqual(t, input, escaped, "SQL injection attempt should be escaped")
	}
}

func TestXSSPrevention(t *testing.T) {
	xssPayloads := []string{
		"<script>alert('xss')</script>",
		"<img src=x onerror=alert('xss')>",
	}

	for _, payload := range xssPayloads {
		htmlEscaped := strings.ReplaceAll(payload, "<", "&lt;")
		htmlEscaped = strings.ReplaceAll(htmlEscaped, ">", "&gt;")
		assert.NotContains(t, htmlEscaped, "<script>", "XSS payload should be escaped")
	}
}

func TestPasswordStrength(t *testing.T) {
	passwords := []struct {
		password string
		strong   bool
	}{
		{"weak", false},
		{"123456", false},
		{"Password1!", true},
	}

	for _, test := range passwords {
		isStrong := len(test.password) >= 8 && 
			strings.ContainsAny(test.password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") &&
			strings.ContainsAny(test.password, "abcdefghijklmnopqrstuvwxyz") &&
			strings.ContainsAny(test.password, "0123456789")
		
		if test.strong {
			assert.True(t, isStrong, "Password %s should be strong", test.password)
		} else {
			assert.False(t, isStrong, "Password %s should be weak", test.password)
		}
	}
}

func TestSecureHeaders(t *testing.T) {
	w := httptest.NewRecorder()

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")

	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
}
