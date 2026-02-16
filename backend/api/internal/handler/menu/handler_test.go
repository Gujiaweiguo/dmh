package menu

import (
	"bytes"
	"context"
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

func setupMenuHandlerTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Menu{}, &model.Role{}, &model.RoleMenu{}, &model.User{}, &model.UserRole{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestMenuHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetMenusHandler(nil))
	assert.NotNil(t, GetMenuHandler(nil))
	assert.NotNil(t, CreateMenuHandler(nil))
	assert.NotNil(t, UpdateMenuHandler(nil))
	assert.NotNil(t, DeleteMenuHandler(nil))
	assert.NotNil(t, ConfigRoleMenusHandler(nil))
	assert.NotNil(t, GetUserMenusHandler(nil))
}

func TestGetMenusHandler_Success(t *testing.T) {
	db := setupMenuHandlerTestDB(t)

	menu := &model.Menu{Name: "Dashboard", Code: "dashboard", Path: "/dashboard", Type: "menu", Platform: "admin", Status: "active"}
	db.Create(menu)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMenusHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/menus", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestCreateMenuHandler_Success(t *testing.T) {
	db := setupMenuHandlerTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := CreateMenuHandler(svcCtx)

	reqBody := types.CreateMenuReq{
		Name:     "New Menu",
		Code:     "new_menu",
		Path:     "/new",
		Type:     "menu",
		Platform: "admin",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/menus", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestCreateMenuHandler_ParseError(t *testing.T) {
	handler := CreateMenuHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/menus", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestUpdateMenuHandler_ParseError(t *testing.T) {
	db := setupMenuHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := UpdateMenuHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/menus/1", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestDeleteMenuHandler_ParseError(t *testing.T) {
	db := setupMenuHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := DeleteMenuHandler(svcCtx)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/menus/invalid", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestConfigRoleMenusHandler_ParseError(t *testing.T) {
	handler := ConfigRoleMenusHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/roles/menus", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetUserMenusHandler_ParseError(t *testing.T) {
	db := setupMenuHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetUserMenusHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/invalid/menus", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetMenuHandler_Success(t *testing.T) {
	db := setupMenuHandlerTestDB(t)

	menu := &model.Menu{Name: "Dashboard", Code: "dashboard", Path: "/dashboard", Type: "menu", Platform: "admin", Status: "active"}
	db.Create(menu)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMenuHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/menus/%d", menu.ID), nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetMenuHandler_NotFound(t *testing.T) {
	db := setupMenuHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMenuHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/menus/99999", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestDeleteMenuHandler_Success(t *testing.T) {
	db := setupMenuHandlerTestDB(t)

	menu := &model.Menu{Name: "Dashboard", Code: "dashboard", Path: "/dashboard", Type: "menu", Platform: "admin", Status: "active"}
	db.Create(menu)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := DeleteMenuHandler(svcCtx)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/menus/%d", menu.ID), nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetMenusHandler_EmptyList(t *testing.T) {
	db := setupMenuHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMenusHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/menus", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetUserMenusHandler_Success(t *testing.T) {
	db := setupMenuHandlerTestDB(t)

	user := &model.User{Username: "testuser", Password: "pass", Phone: "13800138000", Status: "active"}
	db.Create(user)

	role := &model.Role{Name: "admin", Code: "admin"}
	db.Create(role)

	menu := &model.Menu{Name: "Dashboard", Code: "dashboard", Path: "/dashboard", Type: "menu", Platform: "admin", Status: "active"}
	db.Create(menu)

	userRole := &model.UserRole{UserID: user.Id, RoleID: role.ID}
	db.Create(userRole)

	roleMenu := &model.RoleMenu{RoleID: role.ID, MenuID: menu.ID}
	db.Create(roleMenu)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetUserMenusHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%d/menus", user.Id), nil)
	ctx := context.WithValue(req.Context(), "userId", user.Id)
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetMenuHandler_ParseError(t *testing.T) {

	handler := GetMenuHandler(nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/menus/invalid", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestUpdateMenuHandler_Success(t *testing.T) {
	db := setupMenuHandlerTestDB(t)
	// Create a sample menu to update
	menu := &model.Menu{Name: "Dashboard", Code: "dashboard", Path: "/dashboard", Type: "menu", Platform: "admin", Status: "active"}
	db.Create(menu)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := UpdateMenuHandler(svcCtx)

	reqBody := types.UpdateMenuReq{
		Name:     "Dashboard Updated",
		Path:     "/dashboard-updated",
		Type:     "menu",
		Platform: "admin",
		Status:   "active",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/menus/%d", menu.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestConfigRoleMenusHandler_Success(t *testing.T) {
	db := setupMenuHandlerTestDB(t)
	// Seed a role and a menu to associate
	role := &model.Role{Name: "admin", Code: "admin"}
	db.Create(role)
	menu := &model.Menu{Name: "Dashboard", Code: "dashboard", Path: "/dashboard", Type: "menu", Platform: "admin", Status: "active"}
	db.Create(menu)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := ConfigRoleMenusHandler(svcCtx)

	reqBody := types.RoleMenuReq{RoleId: role.ID, MenuIds: []int64{menu.ID}}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/roles/menus", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
