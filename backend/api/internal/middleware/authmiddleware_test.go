package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddlewarePublicPathBypass(t *testing.T) {
	m := NewAuthMiddleware("test-secret")
	called := false
	h := m.Handle(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp := httptest.NewRecorder()
	h(resp, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestAuthMiddlewareProtectedPathUnauthorized(t *testing.T) {
	m := NewAuthMiddleware("test-secret")
	h := m.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	resp := httptest.NewRecorder()
	h(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Token提取失败")
}

func TestAuthMiddlewareValidTokenInjectsContext(t *testing.T) {
	m := NewAuthMiddleware("test-secret")
	token, err := m.GenerateToken(7, "tester", []string{"participant"}, nil)
	require.NoError(t, err)

	var capturedUserID int64
	h := m.Handle(func(w http.ResponseWriter, r *http.Request) {
		uid, uidErr := GetUserIDFromContext(r.Context())
		require.NoError(t, uidErr)
		capturedUserID = uid
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	h(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, int64(7), capturedUserID)
}

func TestGetUserIDFromContextVariants(t *testing.T) {
	ctx1 := context.WithValue(context.Background(), "userId", "12")
	uid1, err1 := GetUserIDFromContext(ctx1)
	assert.NoError(t, err1)
	assert.Equal(t, int64(12), uid1)

	ctx2 := context.WithValue(context.Background(), "userId", 13)
	uid2, err2 := GetUserIDFromContext(ctx2)
	assert.NoError(t, err2)
	assert.Equal(t, int64(13), uid2)
}
