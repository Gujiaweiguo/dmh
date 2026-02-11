package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryRateLimiterAllowAndRemaining(t *testing.T) {
	l := NewMemoryRateLimiter(2, time.Minute)
	uid := "u1"

	assert.True(t, l.Allow(uid))
	assert.True(t, l.Allow(uid))
	assert.False(t, l.Allow(uid))
	assert.Equal(t, 0, l.GetRemaining(uid))
}

func TestRateLimitMiddlewareReturns429(t *testing.T) {
	l := NewMemoryRateLimiter(1, time.Minute)
	mw := RateLimitMiddleware(l)

	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/demo", nil)
	req1.RemoteAddr = "127.0.0.1:1000"
	resp1 := httptest.NewRecorder()
	h.ServeHTTP(resp1, req1)
	assert.Equal(t, http.StatusOK, resp1.Code)

	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/demo", nil)
	req2.RemoteAddr = "127.0.0.1:1000"
	resp2 := httptest.NewRecorder()
	h.ServeHTTP(resp2, req2)
	assert.Equal(t, http.StatusTooManyRequests, resp2.Code)
	assert.Contains(t, resp2.Body.String(), "请求过于频繁")
}
