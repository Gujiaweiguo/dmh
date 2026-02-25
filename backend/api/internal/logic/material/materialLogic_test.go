package material

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MaterialLogicTestSuite struct {
	suite.Suite
	db     *gorm.DB
	svcCtx *svc.ServiceContext
}

func (suite *MaterialLogicTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open("root:Admin168@tcp(127.0.0.1:3306)/dmh_test?charset=utf8mb4&parseTime=true&loc=Local"), &gorm.Config{})
	suite.Require().NoError(err)

	err = db.Exec("DROP TABLE IF EXISTS materials").Error
	suite.Require().NoError(err)

	err = db.AutoMigrate(&model.Material{})
	suite.Require().NoError(err)

	suite.db = db
	suite.svcCtx = &svc.ServiceContext{DB: db}
}

func (suite *MaterialLogicTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func (suite *MaterialLogicTestSuite) SetupTest() {
	suite.db.Exec("TRUNCATE TABLE materials")
}

func (suite *MaterialLogicTestSuite) TestGetMaterialList_Success() {
	ctx := context.Background()

	material1 := &model.Material{Name: "Image1.jpg", Type: "image", Url: "/uploads/1.jpg", CreatedAt: time.Now()}
	material2 := &model.Material{Name: "Text1.txt", Type: "text", Content: "sample content", CreatedAt: time.Now()}
	suite.db.Create(material1)
	suite.db.Create(material2)

	l := NewGetMaterialListLogic(ctx, suite.svcCtx)
	req := &types.GetMaterialListReq{Page: 1, PageSize: 10}

	resp, err := l.GetMaterialList(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(2), resp.Total)
	assert.Len(suite.T(), resp.Materials, 2)
}

func (suite *MaterialLogicTestSuite) TestGetMaterialList_FilterByType() {
	ctx := context.Background()

	suite.db.Create(&model.Material{Name: "Image1.jpg", Type: "image", Url: "/uploads/1.jpg", CreatedAt: time.Now()})
	suite.db.Create(&model.Material{Name: "Image2.jpg", Type: "image", Url: "/uploads/2.jpg", CreatedAt: time.Now()})
	suite.db.Create(&model.Material{Name: "Text1.txt", Type: "text", Content: "content", CreatedAt: time.Now()})

	l := NewGetMaterialListLogic(ctx, suite.svcCtx)
	req := &types.GetMaterialListReq{Page: 1, PageSize: 10, Type: "image"}

	resp, err := l.GetMaterialList(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(2), resp.Total)
	assert.Len(suite.T(), resp.Materials, 2)
	for _, m := range resp.Materials {
		assert.Equal(suite.T(), "image", m.Type)
	}
}

func (suite *MaterialLogicTestSuite) TestGetMaterialList_Pagination() {
	ctx := context.Background()

	for i := 1; i <= 5; i++ {
		suite.db.Create(&model.Material{Name: "Material.jpg", Type: "image", Url: "/uploads/x.jpg", CreatedAt: time.Now()})
	}

	l := NewGetMaterialListLogic(ctx, suite.svcCtx)
	req := &types.GetMaterialListReq{Page: 1, PageSize: 2}

	resp, err := l.GetMaterialList(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(5), resp.Total)
	assert.Len(suite.T(), resp.Materials, 2)
}

func (suite *MaterialLogicTestSuite) TestGetMaterialList_DefaultPagination() {
	ctx := context.Background()

	for i := 1; i <= 25; i++ {
		suite.db.Create(&model.Material{Name: "Material.jpg", Type: "image", Url: "/uploads/x.jpg", CreatedAt: time.Now()})
	}

	l := NewGetMaterialListLogic(ctx, suite.svcCtx)
	req := &types.GetMaterialListReq{}

	resp, err := l.GetMaterialList(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(25), resp.Total)
	assert.Len(suite.T(), resp.Materials, 20)
}

func (suite *MaterialLogicTestSuite) TestGetMaterialList_EmptyResult() {
	ctx := context.Background()

	l := NewGetMaterialListLogic(ctx, suite.svcCtx)
	req := &types.GetMaterialListReq{Page: 1, PageSize: 10}

	resp, err := l.GetMaterialList(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(0), resp.Total)
	assert.Len(suite.T(), resp.Materials, 0)
}

func (suite *MaterialLogicTestSuite) TestDeleteMaterial_Success() {
	ctx := context.Background()

	material := &model.Material{Name: "ToDelete.jpg", Type: "image", Url: "/uploads/del.jpg", CreatedAt: time.Now()}
	suite.db.Create(material)

	l := NewDeleteMaterialLogic(ctx, suite.svcCtx)
	req := &types.DeleteMaterialReq{Id: material.Id}

	resp, err := l.DeleteMaterial(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.True(suite.T(), resp.Success)

	var count int64
	suite.db.Model(&model.Material{}).Where("id = ?", material.Id).Count(&count)
	assert.Equal(suite.T(), int64(0), count)
}

func (suite *MaterialLogicTestSuite) TestDeleteMaterial_NotFound() {
	ctx := context.Background()

	l := NewDeleteMaterialLogic(ctx, suite.svcCtx)
	req := &types.DeleteMaterialReq{Id: 99999}

	resp, err := l.DeleteMaterial(req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "not found")
}

func (suite *MaterialLogicTestSuite) TestGetMaterialList_OrderDesc() {
	ctx := context.Background()

	now := time.Now()
	suite.db.Create(&model.Material{Name: "Old.jpg", Type: "image", Url: "/1.jpg", CreatedAt: now.Add(-time.Hour)})
	suite.db.Create(&model.Material{Name: "New.jpg", Type: "image", Url: "/2.jpg", CreatedAt: now})

	l := NewGetMaterialListLogic(ctx, suite.svcCtx)
	req := &types.GetMaterialListReq{Page: 1, PageSize: 10}

	resp, err := l.GetMaterialList(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "New.jpg", resp.Materials[0].Name)
	assert.Equal(suite.T(), "Old.jpg", resp.Materials[1].Name)
}

// === UploadMaterial Tests ===

func (suite *MaterialLogicTestSuite) TestUploadMaterial_ImageSuccess() {
	ctx := context.Background()

	// Create a test image file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file field
	part, _ := writer.CreateFormFile("file", "test.jpg")
	part.Write([]byte("fake image content"))
	writer.Close()

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/api/material", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	l := NewUploadMaterialLogic(ctx, suite.svcCtx)
	uploadReq := &types.UploadMaterialReq{
		Type: "image",
	}

	resp, err := l.UploadMaterial(uploadReq, req)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "test.jpg", resp.Name)
	assert.Equal(suite.T(), "image", resp.Type)
	assert.NotEmpty(suite.T(), resp.Url)
}

func (suite *MaterialLogicTestSuite) TestUploadMaterial_TextSuccess() {
	ctx := context.Background()

	// Create a test text file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("sample text content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/material", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	l := NewUploadMaterialLogic(ctx, suite.svcCtx)
	uploadReq := &types.UploadMaterialReq{
		Type: "text",
	}

	resp, err := l.UploadMaterial(uploadReq, req)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "test.txt", resp.Name)
	assert.Equal(suite.T(), "text", resp.Type)
	assert.NotEmpty(suite.T(), resp.Id) // Text materials stored in DB
}

func (suite *MaterialLogicTestSuite) TestUploadMaterial_DefaultType() {
	ctx := context.Background()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("file", "default.png")
	part.Write([]byte("image data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/material", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	l := NewUploadMaterialLogic(ctx, suite.svcCtx)
	uploadReq := &types.UploadMaterialReq{
		Type: "", // Empty type should default to "image"
	}

	resp, err := l.UploadMaterial(uploadReq, req)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "image", resp.Type) // Should default to image
}

func (suite *MaterialLogicTestSuite) TestUploadMaterial_NoFile() {
	ctx := context.Background()

	// Create request without file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/material", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	l := NewUploadMaterialLogic(ctx, suite.svcCtx)
	uploadReq := &types.UploadMaterialReq{
		Type: "image",
	}

	resp, err := l.UploadMaterial(uploadReq, req)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "no file uploaded")
}

func TestMaterialLogicTestSuite(t *testing.T) {
	suite.Run(t, new(MaterialLogicTestSuite))
}
