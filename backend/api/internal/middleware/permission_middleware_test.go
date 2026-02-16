package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"dmh/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type permissionTestEnv struct {
	db         *gorm.DB
	middleware *PermissionMiddleware
}

func newPermissionTestEnv(t *testing.T) *permissionTestEnv {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(
		&model.Role{},
		&model.Permission{},
		&model.UserRole{},
		&model.RolePermission{},
	))

	seedPermissionData(t, db)

	sqlDB, err := db.DB()
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, sqlDB.Close())
	})

	return &permissionTestEnv{
		db:         db,
		middleware: NewPermissionMiddleware(sqlDB),
	}
}

func seedPermissionData(t *testing.T, db *gorm.DB) {
	t.Helper()

	roles := []model.Role{
		{ID: 1, Name: "平台管理员", Code: "platform_admin"},
		{ID: 2, Name: "参与者", Code: "participant"},
	}
	permissions := []model.Permission{
		{ID: 1, Name: "订单读取", Code: "orders:read", Resource: "orders", Action: "read"},
		{ID: 2, Name: "品牌读取", Code: "brands:read", Resource: "brands", Action: "read"},
	}
	rolePermissions := []model.RolePermission{
		{RoleID: 1, PermissionID: 1},
		{RoleID: 1, PermissionID: 2},
		{RoleID: 2, PermissionID: 1},
	}
	userRoles := []model.UserRole{
		{UserID: 1, RoleID: 1},
		{UserID: 2, RoleID: 2},
	}

	for _, role := range roles {
		require.NoError(t, db.Create(&role).Error)
	}
	for _, permission := range permissions {
		require.NoError(t, db.Create(&permission).Error)
	}
	for _, rolePermission := range rolePermissions {
		require.NoError(t, db.Create(&rolePermission).Error)
	}
	for _, userRole := range userRoles {
		require.NoError(t, db.Create(&userRole).Error)
	}
}

func TestPermissionMiddleware_CheckPermissionBoundaries(t *testing.T) {
	env := newPermissionTestEnv(t)

	hasPermission, err := env.middleware.checkPermission(context.Background(), 1, "brands:read")
	require.NoError(t, err)
	assert.True(t, hasPermission)

	hasPermission, err = env.middleware.checkPermission(context.Background(), 2, "orders:read")
	require.NoError(t, err)
	assert.True(t, hasPermission)

	hasPermission, err = env.middleware.checkPermission(context.Background(), 2, "brands:read")
	require.NoError(t, err)
	assert.False(t, hasPermission)
}

func TestPermissionMiddleware_ClearUserCacheReloadsPermissions(t *testing.T) {
	env := newPermissionTestEnv(t)

	userPerms, err := env.middleware.getUserPermissions(2)
	require.NoError(t, err)
	assert.False(t, userPerms.Permissions["brands:read"])

	// 动态授予 participant 读品牌权限，验证缓存清除后可见。
	require.NoError(t, env.db.Create(&model.RolePermission{RoleID: 2, PermissionID: 2}).Error)

	userPerms, err = env.middleware.getUserPermissions(2)
	require.NoError(t, err)
	assert.False(t, userPerms.Permissions["brands:read"])

	env.middleware.ClearUserCache(2)
	userPerms, err = env.middleware.getUserPermissions(2)
	require.NoError(t, err)
	assert.True(t, userPerms.Permissions["brands:read"])
}

func TestPermissionMiddleware_HandleRequiresUserContext(t *testing.T) {
	env := newPermissionTestEnv(t)

	handler := env.middleware.Handle(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	resp := httptest.NewRecorder()
	handler(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestPermissionMiddleware_HandleAllowsParticipantOrderRead(t *testing.T) {
	env := newPermissionTestEnv(t)

	handler := env.middleware.Handle(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ctx := context.WithValue(context.Background(), "userId", int64(2))
	ctx = context.WithValue(ctx, "roles", []string{"participant"})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil).WithContext(ctx)
	resp := httptest.NewRecorder()
	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestPermissionMiddleware_HandleBrandIsolationForParticipant(t *testing.T) {
	env := newPermissionTestEnv(t)

	// 先给 participant 授予 brands:read，确保失败来自数据级隔离而不是权限码缺失。
	require.NoError(t, env.db.Create(&model.RolePermission{RoleID: 2, PermissionID: 2}).Error)
	env.middleware.ClearUserCache(2)

	handler := env.middleware.Handle(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ctx := context.WithValue(context.Background(), "userId", int64(2))
	ctx = context.WithValue(ctx, "roles", []string{"participant"})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/brands", nil).WithContext(ctx)
	resp := httptest.NewRecorder()
	handler(resp, req)

	assert.Equal(t, http.StatusForbidden, resp.Code)
	assert.Contains(t, resp.Body.String(), "普通用户无权访问品牌管理接口")
}

func TestPermissionMiddleware_HandlePlatformAdminBypass(t *testing.T) {
	env := newPermissionTestEnv(t)

	handler := env.middleware.Handle(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ctx := context.WithValue(context.Background(), "userId", int64(1))
	ctx = context.WithValue(ctx, "roles", []string{"platform_admin"})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/brands", nil).WithContext(ctx)
	resp := httptest.NewRecorder()
	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestValidatePermission(t *testing.T) {
	assert.NoError(t, ValidatePermission("orders:read"))
	assert.Error(t, ValidatePermission("orders-read"))
	assert.Error(t, ValidatePermission("orders:"))
}

func TestPermissionMiddleware_ExtractResourceAction(t *testing.T) {
	env := newPermissionTestEnv(t)

	tests := []struct {
		name    string
		method  string
		path    string
		wantRes string
		wantAct string
	}{
		{"GET orders", "GET", "/api/v1/orders", "orders", "read"},
		{"POST brands", "POST", "/api/v1/brands", "brands", "create"},
		{"GET campaigns", "GET", "/api/v1/campaigns", "campaigns", "read"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			resource, action := env.middleware.extractResourceAction(req)
			assert.Equal(t, tt.wantRes, resource)
			assert.Equal(t, tt.wantAct, action)
		})
	}
}

func TestIsPublicPath(t *testing.T) {
	env := newPermissionTestEnv(t)

	publicPaths := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/auth/refresh",
		"/health",
		"/ping",
	}

	privatePaths := []string{
		"/api/v1/admin",
		"/api/v1/brands",
		"/api/v1/users",
		"/api/v1/posters/",
	}

	for _, path := range publicPaths {
		t.Run("public_"+path, func(t *testing.T) {
			assert.True(t, env.middleware.isPublicPath(path))
		})
	}

	for _, path := range privatePaths {
		t.Run("private_"+path, func(t *testing.T) {
			assert.False(t, env.middleware.isPublicPath(path))
		})
	}
}

func TestGetUserRolesFromContext(t *testing.T) {
	t.Run("with valid roles slice", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "roles", []string{"admin", "user"})
		roles, err := GetUserRolesFromContext(ctx)
		require.NoError(t, err)
		assert.Equal(t, []string{"admin", "user"}, roles)
	})

	t.Run("with empty roles", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "roles", []string{})
		roles, err := GetUserRolesFromContext(ctx)
		require.NoError(t, err)
		assert.Equal(t, []string{}, roles)
	})

	t.Run("without roles in context", func(t *testing.T) {
		ctx := context.Background()
		_, err := GetUserRolesFromContext(ctx)
		assert.Error(t, err)
	})

	t.Run("with invalid roles type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "roles", "invalid")
		_, err := GetUserRolesFromContext(ctx)
		assert.Error(t, err)
	})
}

func TestExtractBrandIDFromPath(t *testing.T) {
	env := newPermissionTestEnv(t)

	tests := []struct {
		name     string
		path     string
		expected int64
	}{
		{"valid brand id", "/api/v1/brands/123", 123},
		{"brand id with trailing slash", "/api/v1/brands/456/", 456},
		{"brand id with subpath", "/api/v1/brands/789/campaigns", 789},
		{"no brand id in path", "/api/v1/orders", 0},
		{"empty after brands", "/api/v1/brands/", 0},
		{"invalid brand id", "/api/v1/brands/abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := env.middleware.extractBrandIDFromPath(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestClearAllCache(t *testing.T) {
	env := newPermissionTestEnv(t)

	env.middleware.cache.Set(1, &UserPermissions{UserID: 1, Permissions: map[string]bool{"test:read": true}, CachedAt: time.Now()})
	env.middleware.cache.Set(2, &UserPermissions{UserID: 2, Permissions: map[string]bool{"test:write": true}, CachedAt: time.Now()})

	assert.NotNil(t, env.middleware.cache.Get(1))
	assert.NotNil(t, env.middleware.cache.Get(2))

	env.middleware.ClearAllCache()

	assert.Nil(t, env.middleware.cache.Get(1))
	assert.Nil(t, env.middleware.cache.Get(2))
}

func TestPermissionCache_Clear(t *testing.T) {
	cache := NewPermissionCache(time.Minute)

	cache.Set(1, &UserPermissions{UserID: 1, CachedAt: time.Now()})
	cache.Set(2, &UserPermissions{UserID: 2, CachedAt: time.Now()})

	assert.NotNil(t, cache.Get(1))
	assert.NotNil(t, cache.Get(2))

	cache.Clear()

	assert.Nil(t, cache.Get(1))
	assert.Nil(t, cache.Get(2))
}

func TestGetUserBrandIDs(t *testing.T) {
	env := newPermissionTestEnv(t)

	ids, err := env.middleware.getUserBrandIDs(1)
	require.NoError(t, err)
	assert.Equal(t, []int64{}, ids)
}
