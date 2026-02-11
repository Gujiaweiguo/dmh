package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginHandlerParseError(t *testing.T) {
	h := LoginHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	h(resp, req)
	assert.NotEqual(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "unexpected EOF")
}

func TestRegisterHandlerParseError(t *testing.T) {
	h := RegisterHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	h(resp, req)
	assert.NotEqual(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "unexpected EOF")
}
