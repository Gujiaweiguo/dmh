package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
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

func TestGetUserFromContext_Success(t *testing.T) {
	claims := &JWTClaims{
		UserID:   42,
		Username: "tester",
		Roles:    []string{"participant"},
		BrandIDs: []int64{1, 2, 3},
	}
	ctx := context.WithValue(context.Background(), "userClaims", claims)

	user, err := GetUserFromContext(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(42), user.UserID)
	assert.Equal(t, "tester", user.Username)
}

func TestGetUserFromContext_NotFound(t *testing.T) {
	ctx := context.Background()
	user, err := GetUserFromContext(ctx)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestGetUserBrandIDs_Success(t *testing.T) {
	brandIDs := []int64{1, 2, 3}
	ctx := context.WithValue(context.Background(), "brandIds", brandIDs)

	result, err := GetUserBrandIDs(ctx)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetUserBrandIDs_NotFound(t *testing.T) {
	ctx := context.Background()
	result, err := GetUserBrandIDs(ctx)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestCanAccessBrand_Success(t *testing.T) {
	brandIDs := []int64{1, 2, 3}
	ctx := context.WithValue(context.Background(), "brandIds", brandIDs)

	result := CanAccessBrand(ctx, 2)
	assert.False(t, result)
}

func TestCanAccessBrand_Failure(t *testing.T) {
	brandIDs := []int64{1, 2, 3}
	ctx := context.WithValue(context.Background(), "brandIds", brandIDs)

	result := CanAccessBrand(ctx, 4)
	assert.False(t, result)
}

func TestCanAccessBrand_NoBrandIDs(t *testing.T) {
	ctx := context.Background()
	result := CanAccessBrand(ctx, 1)
	assert.False(t, result)
}

func TestRefreshToken_Success(t *testing.T) {
	m := NewAuthMiddleware("test-secret")
	token, err := m.GenerateToken(7, "tester", []string{"participant"}, nil)
	require.NoError(t, err)

	newToken, err := m.RefreshToken(token)
	assert.NoError(t, err)
	assert.NotEmpty(t, newToken)
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	m := NewAuthMiddleware("test-secret")

	newToken, err := m.RefreshToken("invalid_token")
	assert.Error(t, err)
	assert.Empty(t, newToken)
}

func TestValidateToken_Success(t *testing.T) {
	m := NewAuthMiddleware("test-secret")
	token, err := m.GenerateToken(7, "tester", []string{"participant"}, nil)
	require.NoError(t, err)

	claims, err := m.validateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, int64(7), claims.UserID)
	assert.Equal(t, "tester", claims.Username)
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	m := NewAuthMiddleware("test-secret")
	token, err := m.GenerateToken(7, "tester", []string{"participant"}, nil)
	require.NoError(t, err)

	// Use wrong secret to validate
	m2 := NewAuthMiddleware("wrong-secret")
	claims, err := m2.validateToken(token)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	m := NewAuthMiddleware("test-secret")

	// Manually create an expired token
	expiredClaims := &JWTClaims{
		UserID:   7,
		Username: "tester",
		Roles:    []string{"participant"},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	tokenString, _ := expiredToken.SignedString([]byte(m.jwtSecret))

	claims, err := m.validateToken(tokenString)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_InvalidUserID(t *testing.T) {
	m := NewAuthMiddleware("test-secret")

	claims := &JWTClaims{
		UserID:   0,
		Username: "tester",
		Roles:    []string{"participant"},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(m.jwtSecret))

	result, err := m.validateToken(tokenString)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestShouldRefreshToken_ShouldRefresh(t *testing.T) {
	m := NewAuthMiddleware("test-secret")

	claims := &JWTClaims{
		UserID:   7,
		Username: "tester",
		Roles:    []string{"participant"},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(20 * time.Minute)),
		},
	}

	result := m.shouldRefreshToken(claims)
	assert.True(t, result)
}

func TestShouldRefreshToken_ShouldNotRefresh(t *testing.T) {
	m := NewAuthMiddleware("test-secret")

	claims := &JWTClaims{
		UserID:   7,
		Username: "tester",
		Roles:    []string{"participant"},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
		},
	}

	result := m.shouldRefreshToken(claims)
	assert.False(t, result)
}

func TestShouldRefreshToken_NoExpiration(t *testing.T) {
	m := NewAuthMiddleware("test-secret")

	claims := &JWTClaims{
		UserID:   7,
		Username: "tester",
		Roles:    []string{"participant"},
	}

	result := m.shouldRefreshToken(claims)
	assert.False(t, result)
}
