package order

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrderHandlerParseError(t *testing.T) {
	h := CreateOrderHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	h(resp, req)
	assert.NotEqual(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "unexpected EOF")
}

func TestVerifyOrderHandlerParseError(t *testing.T) {
	h := VerifyOrderHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/verify", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	h(resp, req)
	assert.NotEqual(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "unexpected EOF")
}
