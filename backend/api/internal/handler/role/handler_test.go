package role

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRoleHandlerTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Role{}, &model.Permission{}, &model.RolePermission{}, &model.User{}, &model.UserRole{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestRoleHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetRolesHandler(nil))
	assert.NotNil(t, GetPermissionsHandler(nil))
	assert.NotNil(t, ConfigRolePermissionsHandler(nil))
	assert.NotNil(t, GetUserPermissionsHandler(nil))
}

func TestGetRolesHandler_Success(t *testing.T) {
	db := setupRoleHandlerTestDB(t)

	role := &model.Role{Name: "Admin", Code: "admin"}
	db.Create(role)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetRolesHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/roles", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetPermissionsHandler_Success(t *testing.T) {
	db := setupRoleHandlerTestDB(t)

	permission := &model.Permission{Name: "View Users", Code: "users:view", Resource: "users", Action: "read"}
	db.Create(permission)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPermissionsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/permissions", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestConfigRolePermissionsHandler_Success(t *testing.T) {
	db := setupRoleHandlerTestDB(t)
	role := &model.Role{ID: 1, Name: "Platform Admin", Code: "platform_admin"}
	perm1 := &model.Permission{ID: 1, Name: "View Campaign", Code: "campaign:read", Resource: "campaign", Action: "read"}
	perm2 := &model.Permission{ID: 2, Name: "Create Campaign", Code: "campaign:create", Resource: "campaign", Action: "create"}
	db.Create(role)
	db.Create(perm1)
	db.Create(perm2)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := ConfigRolePermissionsHandler(svcCtx)

	body := `{"roleId":1,"permissionIds":[1,2]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/roles/permissions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var commonResp types.CommonResp
	err := json.Unmarshal(resp.Body.Bytes(), &commonResp)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	assert.Equal(t, "Role permissions configured successfully", commonResp.Message)

	var count int64
	db.Model(&model.RolePermission{}).Where("role_id = ?", 1).Count(&count)
	assert.Equal(t, int64(2), count)
}

func TestGetRolesHandler_Success_Content(t *testing.T) {
	db := setupRoleHandlerTestDB(t)
	role := &model.Role{Name: "Admin", Code: "admin"}
	db.Create(role)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetRolesHandler(svcCtx)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/roles", nil)
	resp := httptest.NewRecorder()
	handler(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	var roles []types.RoleResp
	if err := json.Unmarshal(resp.Body.Bytes(), &roles); err != nil {
		t.Fatalf("Failed to unmarshal roles: %v", err)
	}
	assert.NotZero(t, len(roles))
	assert.NotEmpty(t, roles[0].Name)
}

func TestGetPermissionsHandler_Success_Content(t *testing.T) {
	db := setupRoleHandlerTestDB(t)
	permission := &model.Permission{Name: "View Users", Code: "users:view", Resource: "users", Action: "read"}
	db.Create(permission)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPermissionsHandler(svcCtx)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/permissions", nil)
	resp := httptest.NewRecorder()
	handler(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	var perms []types.PermissionResp
	if err := json.Unmarshal(resp.Body.Bytes(), &perms); err != nil {
		t.Fatalf("Failed to unmarshal permissions: %v", err)
	}
	assert.NotZero(t, len(perms))
	assert.NotEmpty(t, perms[0].Name)
}

func TestConfigRolePermissionsHandler_ParseError(t *testing.T) {
	handler := ConfigRolePermissionsHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/roles/permissions", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetUserPermissionsHandler_ParseError(t *testing.T) {
	db := setupRoleHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetUserPermissionsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/invalid/permissions", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetUserPermissionsHandler_Success(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	err = db.AutoMigrate(&model.Role{}, &model.Permission{}, &model.RolePermission{}, &model.User{}, &model.UserRole{}, &model.UserBrand{}, &model.Brand{})
	if err != nil {
		t.Fatalf("Failed to migrate test tables: %v", err)
	}

	user := &model.User{Id: 1, Username: "testuser", Phone: "13800138000"}
	role := &model.Role{ID: 1, Name: "平台管理员", Code: "platform_admin"}
	permission := &model.Permission{ID: 1, Name: "查看活动", Code: "campaign:read", Resource: "campaign", Action: "read"}
	brand := &model.Brand{Id: 1, Name: "Test Brand"}
	db.Create(user)
	db.Create(role)
	db.Create(permission)
	db.Create(brand)

	userRole := &model.UserRole{UserID: 1, RoleID: 1}
	db.Create(userRole)
	rolePerm := &model.RolePermission{RoleID: 1, PermissionID: 1}
	db.Create(rolePerm)
	userBrand := &model.UserBrand{UserId: 1, BrandId: 1}
	db.Create(userBrand)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetUserPermissionsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/1/permissions", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result types.UserPermissionResp
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, int64(1), result.UserId)
	assert.Contains(t, result.Roles, "platform_admin")
	assert.Contains(t, result.Permissions, "campaign:read")
	assert.Contains(t, result.BrandIds, int64(1))
}

func TestGetUserPermissionsHandler_NoPermissions_Success(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Role{}, &model.Permission{}, &model.RolePermission{}, &model.UserRole{}, &model.UserBrand{}, &model.Brand{}); err != nil {
		t.Fatalf("Failed to migrate test tables: %v", err)
	}

	user := &model.User{Id: 2, Username: "emptyuser", Phone: "123"}
	db.Create(user)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetUserPermissionsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/2/permissions", nil)
	resp := httptest.NewRecorder()
	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var result types.UserPermissionResp
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	assert.Equal(t, int64(2), result.UserId)
	assert.Empty(t, result.Roles)
	assert.Empty(t, result.Permissions)
	assert.Empty(t, result.BrandIds)
}

func TestGetUserPermissionsHandler_Success_MultiRolesAndPerms(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	if err := db.AutoMigrate(&model.Role{}, &model.Permission{}, &model.RolePermission{}, &model.User{}, &model.UserRole{}, &model.UserBrand{}, &model.Brand{}); err != nil {
		t.Fatalf("Failed to migrate test tables: %v", err)
	}

	role1 := &model.Role{ID: 1, Name: "Role One", Code: "role_one"}
	role2 := &model.Role{ID: 2, Name: "Role Two", Code: "role_two"}
	perm1 := &model.Permission{ID: 1, Name: "View Campaign", Code: "campaign:read", Resource: "campaign", Action: "read"}
	perm2 := &model.Permission{ID: 2, Name: "Create Campaign", Code: "campaign:create", Resource: "campaign", Action: "create"}
	db.Create(role1)
	db.Create(role2)
	db.Create(perm1)
	db.Create(perm2)

	user := &model.User{Id: 3, Username: "multirole", Phone: "123"}
	db.Create(user)
	ur1 := &model.UserRole{UserID: 3, RoleID: 1}
	ur2 := &model.UserRole{UserID: 3, RoleID: 2}
	db.Create(ur1)
	db.Create(ur2)
	pr1 := &model.RolePermission{RoleID: 1, PermissionID: 1}
	pr2 := &model.RolePermission{RoleID: 2, PermissionID: 2}
	db.Create(pr1)
	db.Create(pr2)
	b1 := &model.UserBrand{UserId: 3, BrandId: 1}
	b2 := &model.UserBrand{UserId: 3, BrandId: 2}
	db.Create(b1)
	db.Create(b2)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetUserPermissionsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/3/permissions", nil)
	resp := httptest.NewRecorder()
	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result types.UserPermissionResp
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	assert.Equal(t, int64(3), result.UserId)

	assert.Contains(t, result.Roles, "role_one")
	assert.Contains(t, result.Roles, "role_two")

	assert.Contains(t, result.Permissions, "campaign:read")
	assert.Contains(t, result.Permissions, "campaign:create")

	assert.Contains(t, result.BrandIds, int64(1))
	assert.Contains(t, result.BrandIds, int64(2))
}
