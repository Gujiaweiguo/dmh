package member

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
	"gorm.io/gorm"
)

func setupMemberHandlerTestDB(t *testing.T) *gorm.DB {
	db := testutil.SetupGormTestDB(t)

	err := db.AutoMigrate(&model.Member{}, &model.MemberProfile{}, &model.MemberTag{}, &model.Brand{}, &model.Campaign{}, &model.Order{})
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

	member := &model.Member{UnionID: testutil.GenUniqueUnionID(), Nickname: "Test Member", Phone: testutil.GenUniquePhone(), Status: "active"}
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

	member := &model.Member{UnionID: testutil.GenUniqueUnionID(), Nickname: "Test Member", Phone: testutil.GenUniquePhone(), Status: "active"}
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

	member := &model.Member{UnionID: testutil.GenUniqueUnionID(), Nickname: "Test Member", Phone: testutil.GenUniquePhone(), Status: "active"}
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

	member := &model.Member{UnionID: testutil.GenUniqueUnionID(), Nickname: "Test Member", Phone: testutil.GenUniquePhone(), Status: "active"}
	db.Create(member)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMembersHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/members?page=1&pageSize=10&phone=%s", member.Phone), nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestUpdateMemberHandler_Success(t *testing.T) {
	db := setupMemberHandlerTestDB(t)

	member := &model.Member{UnionID: testutil.GenUniqueUnionID(), Nickname: "Test Member", Phone: testutil.GenUniquePhone(), Status: "active"}
	db.Create(member)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := UpdateMemberHandler(svcCtx)

	reqBody := types.UpdateMemberReq{
		Nickname: "Updated Name",
		Avatar:   "https://example.com/avatar.png",
		Gender:   1,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/members/%d", member.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestUpdateMemberHandler_InvalidPath(t *testing.T) {
	db := setupMemberHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := UpdateMemberHandler(svcCtx)

	reqBody := types.UpdateMemberReq{
		Nickname: "Updated Name",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/members/invalid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestUpdateMemberStatusHandler_Success(t *testing.T) {
	db := setupMemberHandlerTestDB(t)

	member := &model.Member{UnionID: testutil.GenUniqueUnionID(), Nickname: "Test Member", Phone: testutil.GenUniquePhone(), Status: "active"}
	db.Create(member)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := UpdateMemberStatusHandler(svcCtx)

	reqBody := types.UpdateMemberStatusReq{
		Status: "inactive",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/members/%d/status", member.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestUpdateMemberStatusHandler_InvalidPath(t *testing.T) {
	db := setupMemberHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := UpdateMemberStatusHandler(svcCtx)

	reqBody := types.UpdateMemberStatusReq{
		Status: "inactive",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/members/invalid/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestParseMemberIDFromPath_Success(t *testing.T) {
	id, err := parseMemberIDFromPath("/api/v1/members/123")
	assert.NoError(t, err)
	assert.Equal(t, int64(123), id)
}

func TestParseMemberIDFromPath_Invalid(t *testing.T) {
	id, err := parseMemberIDFromPath("/api/v1/members/invalid")
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
}
