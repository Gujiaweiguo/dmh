package admin

import (
	"bytes"
	"dmh/api/internal/handler/testutil"
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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupAdminHandlerTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.SetupGormTestDB(t)

	err := db.AutoMigrate(&model.User{}, &model.Role{}, &model.UserRole{}, &model.UserBrand{}, &model.Brand{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestAdminHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetUsersHandler(nil))
	assert.NotNil(t, GetUserHandler(nil))
	assert.NotNil(t, CreateUserHandler(nil))
	assert.NotNil(t, UpdateUserHandler(nil))
	assert.NotNil(t, DeleteUserHandler(nil))
	assert.NotNil(t, ResetUserPasswordHandler(nil))
	assert.NotNil(t, ManageBrandAdminRelationHandler(nil))
}

func TestGetUsersHandler_Success(t *testing.T) {
	db := setupAdminHandlerTestDB(t)

	user := &model.User{Username: testutil.GenUniqueUsername("admin"), Password: "pass", Phone: testutil.GenUniquePhone(), Status: "active", Role: "platform_admin"}
	db.Create(user)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetUsersHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestCreateUserHandler_Success(t *testing.T) {
	db := setupAdminHandlerTestDB(t)

	role := &model.Role{Name: "Participant", Code: "participant"}
	db.Create(role)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := CreateUserHandler(svcCtx)

	reqBody := types.AdminCreateUserReq{
		Username: "newuser",
		Password: "password123",
		Phone:    "13800138111",
		Role:     "participant",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestCreateUserHandler_ParseError(t *testing.T) {
	handler := CreateUserHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetUserHandler_ParseError(t *testing.T) {
	db := setupAdminHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetUserHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users/invalid", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestUpdateUserHandler_ParseError(t *testing.T) {
	db := setupAdminHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := UpdateUserHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/users/1", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestDeleteUserHandler_ParseError(t *testing.T) {
	db := setupAdminHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := DeleteUserHandler(svcCtx)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/users/invalid", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestResetUserPasswordHandler_ParseError(t *testing.T) {
	db := setupAdminHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := ResetUserPasswordHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users/1/reset-password", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestManageBrandAdminRelationHandler_ParseError(t *testing.T) {
	db := setupAdminHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := ManageBrandAdminRelationHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users/brand-relations", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestManageBrandAdminRelationHandler_Success(t *testing.T) {
	db := setupAdminHandlerTestDB(t)
	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	user := &model.User{Username: testutil.GenUniqueUsername("admin"), Password: "pass", Phone: testutil.GenUniquePhone(), Status: "active", Role: "platform_admin"}
	db.Create(user)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := ManageBrandAdminRelationHandler(svcCtx)

	reqBody := types.BrandAdminRelationReq{
		UserId:   user.Id,
		BrandIds: []int64{brand.Id},
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users/brand-relations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestResetUserPasswordHandler_Success(t *testing.T) {
	db := setupAdminHandlerTestDB(t)
	user := &model.User{Username: testutil.GenUniqueUsername("admin"), Password: "pass", Phone: testutil.GenUniquePhone(), Status: "active", Role: "platform_admin"}
	db.Create(user)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := ResetUserPasswordHandler(svcCtx)

	reqBody := types.AdminResetPasswordReq{
		NewPassword: "newpassword123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/admin/users/%d/reset-password", user.Id), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestUpdateUserHandler_Success(t *testing.T) {
	db := setupAdminHandlerTestDB(t)
	user := &model.User{Username: testutil.GenUniqueUsername("admin"), Password: "pass", Phone: testutil.GenUniquePhone(), Status: "active", Role: "platform_admin"}
	db.Create(user)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := UpdateUserHandler(svcCtx)

	reqBody := types.AdminUpdateUserReq{
		RealName: "Admin User",
		Status:   "active",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/admin/users/%d", user.Id), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestUpdateUserHandler_UsesPathID(t *testing.T) {
	db := setupAdminHandlerTestDB(t)
	user1 := &model.User{Username: "u1", Password: "pass", Phone: "13800138101", RealName: "用户1", Status: "active"}
	user2 := &model.User{Username: "u2", Password: "pass", Phone: "13800138102", RealName: "用户2", Status: "active"}
	db.Create(user1)
	db.Create(user2)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := UpdateUserHandler(svcCtx)

	reqBody := types.AdminUpdateUserReq{RealName: "被修改"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/admin/users/%d", user2.Id), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var got1, got2 model.User
	db.First(&got1, user1.Id)
	db.First(&got2, user2.Id)
	assert.Equal(t, "用户1", got1.RealName)
	assert.Equal(t, "被修改", got2.RealName)
}

func TestDeleteUserHandler_UsesPathID(t *testing.T) {
	db := setupAdminHandlerTestDB(t)
	user1 := &model.User{Username: testutil.GenUniqueUsername("d1"), Password: "pass", Phone: testutil.GenUniquePhone(), Status: "active"}
	user2 := &model.User{Username: "d2", Password: "pass", Phone: "13800138112", Status: "active"}
	db.Create(user1)
	db.Create(user2)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := DeleteUserHandler(svcCtx)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/admin/users/%d", user2.Id), nil)
	resp := httptest.NewRecorder()

	handler(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var got1, got2 model.User
	err1 := db.First(&got1, user1.Id).Error
	err2 := db.First(&got2, user2.Id).Error
	assert.NoError(t, err1)
	assert.Error(t, err2)
}

func TestResetUserPasswordHandler_UsesPathID(t *testing.T) {
	db := setupAdminHandlerTestDB(t)
	old1, _ := bcrypt.GenerateFromPassword([]byte("oldpass1"), bcrypt.DefaultCost)
	old2, _ := bcrypt.GenerateFromPassword([]byte("oldpass2"), bcrypt.DefaultCost)
	user1 := &model.User{Username: "r1", Password: string(old1), Phone: "13800138121", Status: "active"}
	user2 := &model.User{Username: "r2", Password: string(old2), Phone: "13800138122", Status: "active"}
	db.Create(user1)
	db.Create(user2)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := ResetUserPasswordHandler(svcCtx)

	reqBody := types.AdminResetPasswordReq{NewPassword: "new-password-2"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/admin/users/%d/reset-password", user2.Id), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var got1, got2 model.User
	db.First(&got1, user1.Id)
	db.First(&got2, user2.Id)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(got1.Password), []byte("oldpass1")))
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(got2.Password), []byte("new-password-2")))
}
