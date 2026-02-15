package member

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dmh/api/internal/svc"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupMemberHandlerTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Member{}, &model.MemberProfile{}, &model.MemberTag{}, &model.Brand{}, &model.Campaign{}, &model.Order{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestMemberHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetMembersHandler(nil))
	assert.NotNil(t, GetMemberHandler(nil))
	assert.NotNil(t, GetMemberProfileHandler(nil))
	assert.NotNil(t, UpdateMemberHandler(nil))
	assert.NotNil(t, UpdateMemberStatusHandler(nil))
}

func TestGetMembersHandler_Success(t *testing.T) {
	db := setupMemberHandlerTestDB(t)

	member := &model.Member{UnionID: "union123", Nickname: "Test Member", Phone: "13800138000", Status: "active"}
	db.Create(member)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMembersHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/members?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetMemberHandler_ParseError(t *testing.T) {
	db := setupMemberHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMemberHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/members/invalid", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetMemberProfileHandler_ParseError(t *testing.T) {
	db := setupMemberHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMemberProfileHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/members/invalid/profile", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestUpdateMemberHandler_ParseError(t *testing.T) {
	db := setupMemberHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := UpdateMemberHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/members/1", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestUpdateMemberStatusHandler_ParseError(t *testing.T) {
	db := setupMemberHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := UpdateMemberStatusHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/members/1/status", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetMemberHandler_Success(t *testing.T) {
	db := setupMemberHandlerTestDB(t)

	member := &model.Member{UnionID: "union123", Nickname: "Test Member", Phone: "13800138000", Status: "active"}
	db.Create(member)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMemberHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/members/%d", member.ID), nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetMemberHandler_NotFound(t *testing.T) {
	db := setupMemberHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMemberHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/members/99999", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetMemberProfileHandler_Success(t *testing.T) {
	db := setupMemberHandlerTestDB(t)

	member := &model.Member{UnionID: "union123", Nickname: "Test Member", Phone: "13800138000", Status: "active"}
	db.Create(member)

	profile := &model.MemberProfile{MemberID: member.ID, TotalOrders: 5, TotalPayment: 100.00}
	db.Create(profile)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMemberProfileHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/members/%d/profile", member.ID), nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetMembersHandler_EmptyList(t *testing.T) {
	db := setupMemberHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMembersHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/members?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetMembersHandler_WithFilters(t *testing.T) {
	db := setupMemberHandlerTestDB(t)

	member := &model.Member{UnionID: "union123", Nickname: "Test Member", Phone: "13800138000", Status: "active"}
	db.Create(member)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMembersHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/members?page=1&pageSize=10&phone=13800138000", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
