package middleware

import (
	"context"
	"fmt"
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

type mockRedisClient struct {
	counts  map[string]int64
	expires map[string]int
	ttls    map[string]int
}

func newMockRedisClient() *mockRedisClient {
	return &mockRedisClient{
		counts:  make(map[string]int64),
		expires: make(map[string]int),
		ttls:    make(map[string]int),
	}
}

func (m *mockRedisClient) Incr(ctx context.Context, key string) (int64, error) {
	m.counts[key]++
	return m.counts[key], nil
}

func (m *mockRedisClient) Get(ctx context.Context, key string) (string, error) {
	return fmt.Sprintf("%d", m.counts[key]), nil
}

func (m *mockRedisClient) Expire(ctx context.Context, key string, seconds int) error {
	m.expires[key] = seconds
	return nil
}

func (m *mockRedisClient) TTL(ctx context.Context, key string) (string, error) {
	return fmt.Sprintf("%d", m.ttls[key]), nil
}

func TestRedisRateLimiter_Allow(t *testing.T) {
	mock := newMockRedisClient()
	limiter := NewRedisRateLimiter(mock, "test", 2, time.Minute)

	assert.True(t, limiter.Allow("user1"))
	assert.True(t, limiter.Allow("user1"))
	assert.False(t, limiter.Allow("user1"))
}

func TestRedisRateLimiter_GetRemaining(t *testing.T) {
	mock := newMockRedisClient()
	limiter := NewRedisRateLimiter(mock, "test", 5, time.Minute)

	remaining := limiter.GetRemaining("user1")
	assert.Equal(t, 5, remaining)

	limiter.Allow("user1")
	remaining = limiter.GetRemaining("user1")
	assert.Equal(t, 4, remaining)
}

func TestRedisRateLimiter_GetResetTime(t *testing.T) {
	mock := newMockRedisClient()
	mock.ttls["test:requests:expire:user1"] = 60
	limiter := NewRedisRateLimiter(mock, "test", 2, time.Minute)

	resetTime := limiter.GetResetTime("user1")
	expectedTime := time.Now().Add(60 * time.Second)
	assert.True(t, resetTime.Before(expectedTime.Add(time.Second)))
}

func TestMemoryRateLimiter_GetResetTime(t *testing.T) {
	l := NewMemoryRateLimiter(2, time.Minute)
	l.Allow("user1")

	resetTime := l.GetResetTime("user1")
	now := time.Now()
	assert.True(t, resetTime.After(now) || resetTime.Equal(now))
}

func TestMemoryRateLimiter_GetRemaining_NoUser(t *testing.T) {
	l := NewMemoryRateLimiter(5, time.Minute)

	remaining := l.GetRemaining("nonexistent")
	assert.Equal(t, 5, remaining)
}

func TestMemoryRateLimiter_GetResetTime_NoUser(t *testing.T) {
	l := NewMemoryRateLimiter(2, time.Minute)

	resetTime := l.GetResetTime("nonexistent")
	now := time.Now()
	assert.True(t, resetTime.After(now) || resetTime.Equal(now) || resetTime.Before(now.Add(time.Second)))
}
