package material

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"dmh/api/internal/handler/testutil"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupMaterialHandlerTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.SetupGormTestDB(t)

	err := db.AutoMigrate(&model.Material{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	testutil.ClearTables(db, "materials")

	return db
}

func createTestMaterial(t *testing.T, db *gorm.DB, name string, materialType string) *model.Material {
	t.Helper()
	m := &model.Material{
		Name:      name + fmt.Sprintf("_%d", time.Now().UnixNano()),
		Type:      materialType,
		Url:       "/uploads/test.png",
		CreatedAt: time.Now(),
	}
	if err := db.Create(m).Error; err != nil {
		t.Fatalf("Failed to create test material: %v", err)
	}
	return m
}

func TestGetMaterialListHandler(t *testing.T) {
	db := setupMaterialHandlerTestDB(t)
	_ = createTestMaterial(t, db, "test_image", "image")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMaterialListHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/material/list?page=1&pageSize=10", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp types.MaterialListResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Materials, 1)
	assert.Equal(t, "image", resp.Materials[0].Type)
}

func TestGetMaterialListHandlerWithFilters(t *testing.T) {
	db := setupMaterialHandlerTestDB(t)
	_ = createTestMaterial(t, db, "test_image", "image")
	_ = createTestMaterial(t, db, "test_text", "text")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMaterialListHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/material/list?type=image", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp types.MaterialListResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Materials, 1)
	assert.Equal(t, "image", resp.Materials[0].Type)
}

func TestGetMaterialListHandlerEmpty(t *testing.T) {
	db := setupMaterialHandlerTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMaterialListHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/material/list?page=1&pageSize=10", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp types.MaterialListResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, int64(0), resp.Total)
	assert.Len(t, resp.Materials, 0)
}

func TestDeleteMaterialHandler(t *testing.T) {
	db := setupMaterialHandlerTestDB(t)
	material := createTestMaterial(t, db, "test_to_delete", "image")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := DeleteMaterialHandler(svcCtx)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/material/delete/%d", material.Id), nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp types.DeleteMaterialResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.True(t, resp.Success)

	var count int64
	db.Model(&model.Material{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteMaterialHandlerNotFound(t *testing.T) {
	db := setupMaterialHandlerTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := DeleteMaterialHandler(svcCtx)

	req := httptest.NewRequest(http.MethodDelete, "/material/delete/99999", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.NotEqual(t, http.StatusOK, rec.Code)
}

func TestUploadMaterialHandlerConstruct(t *testing.T) {
	handler := UploadMaterialHandler(nil)
	assert.NotNil(t, handler)
}

func TestUploadMaterialHandler(t *testing.T) {
	db := setupMaterialHandlerTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := UploadMaterialHandler(svcCtx)

	body := &bytes.Buffer{}
	writer := createMultipartFormData(t, body, "file", "test.png", []byte("fake image content"))
	req := httptest.NewRequest(http.MethodPost, "/material/upload?type=image", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp types.UploadMaterialResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.NotZero(t, resp.Id)
	assert.Equal(t, "test.png", resp.Name)
	assert.Equal(t, "image", resp.Type)
	assert.NotEmpty(t, resp.Url)
}

func TestUploadMaterialHandlerNoFile(t *testing.T) {
	db := setupMaterialHandlerTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := UploadMaterialHandler(svcCtx)

	// Request without multipart form data
	body := &bytes.Buffer{}
	req := httptest.NewRequest(http.MethodPost, "/material/upload?type=image", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.NotEqual(t, http.StatusOK, rec.Code)
}

func TestUploadMaterialHandlerNoType(t *testing.T) {
	db := setupMaterialHandlerTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := UploadMaterialHandler(svcCtx)

	body := &bytes.Buffer{}
	writer := createMultipartFormData(t, body, "file", "test.txt", []byte("text content"))
	req := httptest.NewRequest(http.MethodPost, "/material/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp types.UploadMaterialResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, "image", resp.Type) // Default type is "image" when not specified
}

func TestUploadMaterialHandlerTextType(t *testing.T) {
	db := setupMaterialHandlerTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := UploadMaterialHandler(svcCtx)

	body := &bytes.Buffer{}
	writer := createMultipartFormData(t, body, "file", "test.txt", []byte("text content"))
	req := httptest.NewRequest(http.MethodPost, "/material/upload?type=text", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp types.UploadMaterialResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, "text", resp.Type)
	assert.Equal(t, "test.txt", resp.Name)
}

func TestGetMaterialListHandlerPagination(t *testing.T) {
	db := setupMaterialHandlerTestDB(t)

	// Create 25 materials
	for i := 0; i < 25; i++ {
		_ = createTestMaterial(t, db, fmt.Sprintf("test_%d", i), "image")
	}

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMaterialListHandler(svcCtx)

	// Test first page (page 1, pageSize 10)
	req := httptest.NewRequest(http.MethodGet, "/material/list?page=1&pageSize=10", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp types.MaterialListResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, int64(25), resp.Total)
	assert.Len(t, resp.Materials, 10)

	// Test second page
	req = httptest.NewRequest(http.MethodGet, "/material/list?page=2&pageSize=10", nil)
	rec = httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, int64(25), resp.Total)
	assert.Len(t, resp.Materials, 10)
}

func TestGetMaterialListHandlerDefaultPagination(t *testing.T) {
	db := setupMaterialHandlerTestDB(t)
	_ = createTestMaterial(t, db, "test_image", "image")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMaterialListHandler(svcCtx)

	// Request without page/pageSize parameters
	req := httptest.NewRequest(http.MethodGet, "/material/list", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp types.MaterialListResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Materials, 1)
}

func TestGetMaterialListHandlerInvalidPage(t *testing.T) {
	db := setupMaterialHandlerTestDB(t)
	_ = createTestMaterial(t, db, "test_image", "image")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMaterialListHandler(svcCtx)

	// Request with negative page
	req := httptest.NewRequest(http.MethodGet, "/material/list?page=-1&pageSize=10", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp types.MaterialListResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Should default to page 1
	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Materials, 1)
}

func TestGetMaterialListHandlerInvalidPageSize(t *testing.T) {
	db := setupMaterialHandlerTestDB(t)
	_ = createTestMaterial(t, db, "test_image", "image")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMaterialListHandler(svcCtx)

	// Request with 0 pageSize
	req := httptest.NewRequest(http.MethodGet, "/material/list?page=1&pageSize=0", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp types.MaterialListResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Should default to pageSize 20
	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Materials, 1)
}

func TestDeleteMaterialHandlerInvalidID(t *testing.T) {
	db := setupMaterialHandlerTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := DeleteMaterialHandler(svcCtx)

	// Request with non-numeric ID
	req := httptest.NewRequest(http.MethodDelete, "/material/delete/invalid", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.NotEqual(t, http.StatusOK, rec.Code)
}

func TestDeleteMaterialHandlerZeroID(t *testing.T) {
	db := setupMaterialHandlerTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := DeleteMaterialHandler(svcCtx)

	// Request with ID = 0
	req := httptest.NewRequest(http.MethodDelete, "/material/delete/0", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.NotEqual(t, http.StatusOK, rec.Code)
}


func createMultipartFormData(t *testing.T, body *bytes.Buffer, fieldName, filename string, content []byte) *multipart.Writer {
	t.Helper()
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	_, err = part.Write(content)
	if err != nil {
		t.Fatalf("Failed to write content: %v", err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}
	return writer
}
