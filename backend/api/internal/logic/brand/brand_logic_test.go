package brand

import (
	"context"
	"strconv"
	"testing"

	"dmh/api/internal/handler/testutil"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupBrandTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.SetupGormTestDB(t)

	err := db.AutoMigrate(&model.Brand{}, &model.BrandAsset{}, &model.Campaign{}, &model.Order{}, &model.Reward{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	testutil.ClearTables(db, "rewards", "orders", "campaigns", "brand_assets", "brands")

	return db
}

func createTestBrand(t *testing.T, db *gorm.DB, name, status string) *model.Brand {
	brand := &model.Brand{
		Name:        name,
		Description: "Test brand description",
		Status:      status,
	}
	if err := db.Create(brand).Error; err != nil {
		t.Fatalf("Failed to create test brand: %v", err)
	}
	return brand
}

func TestGetBrandsLogic_GetBrands_Success(t *testing.T) {
	db := setupBrandTestDB(t)

	createTestBrand(t, db, "Brand 1", "active")
	createTestBrand(t, db, "Brand 2", "active")
	createTestBrand(t, db, "Brand 3", "disabled")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandsLogic(ctx, svcCtx)

	resp, err := logic.GetBrands()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(3), resp.Total)
	assert.Len(t, resp.Brands, 3)
}

func TestGetBrandsLogic_ReturnsCorrectData(t *testing.T) {
	db := setupBrandTestDB(t)

	brand := createTestBrand(t, db, "Test Brand", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandsLogic(ctx, svcCtx)

	resp, err := logic.GetBrands()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Brands, 1)
	assert.Equal(t, brand.Id, resp.Brands[0].Id)
	assert.Equal(t, brand.Name, resp.Brands[0].Name)
	assert.Equal(t, brand.Status, resp.Brands[0].Status)
}

func TestGetBrandsLogic_EmptyResult(t *testing.T) {
	db := setupBrandTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandsLogic(ctx, svcCtx)

	resp, err := logic.GetBrands()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(0), resp.Total)
	assert.Len(t, resp.Brands, 0)
}

func TestUpdateBrandAssetLogic_InvalidParams(t *testing.T) {
	db := setupBrandTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateBrandAssetLogic(ctx, svcCtx)

	req := &types.UpdateBrandAssetReq{
		BrandId: 0,
		Id:      1,
	}

	resp, err := logic.UpdateBrandAsset(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "参数无效")
}

func TestUpdateBrandAssetLogic_AssetNotFound(t *testing.T) {
	db := setupBrandTestDB(t)

	brand := createTestBrand(t, db, "Test Brand", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateBrandAssetLogic(ctx, svcCtx)

	req := &types.UpdateBrandAssetReq{
		BrandId: brand.Id,
		Id:      999,
		Type:    "image",
	}

	resp, err := logic.UpdateBrandAsset(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "素材不存在")
}

func TestUpdateBrandAssetLogic_NoUpdates(t *testing.T) {
	db := setupBrandTestDB(t)

	brand := createTestBrand(t, db, "Test Brand", "active")
	asset := &model.BrandAsset{
		BrandID:     brand.Id,
		Name:        "Test Asset",
		Type:        "image",
		FileUrl:     "https://example.com/image.png",
		Description: "Test description",
	}
	db.Create(asset)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateBrandAssetLogic(ctx, svcCtx)

	req := &types.UpdateBrandAssetReq{
		BrandId: brand.Id,
		Id:      asset.ID,
	}

	resp, err := logic.UpdateBrandAsset(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, asset.ID, resp.Id)
	assert.Equal(t, "Test Asset", resp.Name)
}

func TestUpdateBrandAssetLogic_Success(t *testing.T) {
	db := setupBrandTestDB(t)

	brand := createTestBrand(t, db, "Test Brand", "active")
	asset := &model.BrandAsset{
		BrandID:     brand.Id,
		Name:        "Test Asset",
		Type:        "image",
		FileUrl:     "https://example.com/image.png",
		Description: "Test description",
	}
	db.Create(asset)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateBrandAssetLogic(ctx, svcCtx)

	req := &types.UpdateBrandAssetReq{
		BrandId:     brand.Id,
		Id:          asset.ID,
		Name:        "Updated Asset",
		Type:        "video",
		Category:    "marketing",
		Tags:        "tag1,tag2",
		FileUrl:     "https://example.com/video.mp4",
		FileSize:    1024,
		Description: "Updated description",
	}

	resp, err := logic.UpdateBrandAsset(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Updated Asset", resp.Name)
	assert.Equal(t, "video", resp.Type)
	assert.Equal(t, "marketing", resp.Category)
}

func TestDeleteBrandAssetLogic_Success(t *testing.T) {
	db := setupBrandTestDB(t)

	brand := createTestBrand(t, db, "Test Brand", "active")
	asset := &model.BrandAsset{
		BrandID: brand.Id,
		Name:    "Test Asset",
		Type:    "image",
		FileUrl: "https://example.com/image.png",
	}
	db.Create(asset)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewDeleteBrandAssetLogic(ctx, svcCtx)

	req := &types.DeleteBrandAssetReq{
		BrandId: brand.Id,
		Id:      asset.ID,
	}

	resp, err := logic.DeleteBrandAsset(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	var count int64
	db.Model(&model.BrandAsset{}).Where("id = ?", asset.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteBrandAssetLogic_NotFound(t *testing.T) {
	db := setupBrandTestDB(t)

	brand := createTestBrand(t, db, "Test Brand", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewDeleteBrandAssetLogic(ctx, svcCtx)

	req := &types.DeleteBrandAssetReq{
		BrandId: brand.Id,
		Id:      999,
	}

	resp, err := logic.DeleteBrandAsset(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetBrandLogic_Success(t *testing.T) {
	db := setupBrandTestDB(t)

	brand := createTestBrand(t, db, "Test Brand", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandLogic(ctx, svcCtx)

	req := &types.GetBrandReq{Id: brand.Id}

	resp, err := logic.GetBrand(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, brand.Name, resp.Name)
	assert.Equal(t, brand.Status, resp.Status)
}

func TestGetBrandLogic_NotFound(t *testing.T) {
	db := setupBrandTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandLogic(ctx, svcCtx)

	req := &types.GetBrandReq{Id: 999}

	resp, err := logic.GetBrand(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetBrandLogic_InvalidId(t *testing.T) {
	db := setupBrandTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandLogic(ctx, svcCtx)

	req := &types.GetBrandReq{Id: 0}

	resp, err := logic.GetBrand(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "品牌ID无效")
}

func TestCreateBrandLogic_Success(t *testing.T) {
	db := setupBrandTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewCreateBrandLogic(ctx, svcCtx)

	req := &types.CreateBrandReq{
		Name:        "New Brand",
		Description: "New brand description",
		Logo:        "https://example.com/logo.png",
	}

	resp, err := logic.CreateBrand(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "New Brand", resp.Name)
	assert.Greater(t, resp.Id, int64(0))
}

func TestCreateBrandLogic_EmptyName(t *testing.T) {
	db := setupBrandTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewCreateBrandLogic(ctx, svcCtx)

	req := &types.CreateBrandReq{
		Name:        "",
		Description: "New brand description",
	}

	resp, err := logic.CreateBrand(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "品牌名称不能为空")
}

func TestCreateBrandLogic_WhitespaceName(t *testing.T) {
	db := setupBrandTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewCreateBrandLogic(ctx, svcCtx)

	req := &types.CreateBrandReq{
		Name:        "   ",
		Description: "New brand description",
	}

	resp, err := logic.CreateBrand(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "品牌名称不能为空")
}

func TestCreateBrandAssetLogic_Success(t *testing.T) {
	db := setupBrandTestDB(t)

	brand := createTestBrand(t, db, "Test Brand", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewCreateBrandAssetLogic(ctx, svcCtx)

	req := &types.BrandAssetReq{
		BrandId:  strconv.FormatInt(brand.Id, 10),
		Type:     "image",
		FileUrl:  "https://example.com/image.png",
		Name:     "Test Asset",
		Category: "marketing",
	}

	resp, err := logic.CreateBrandAsset(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "image", resp.Type)
	assert.Equal(t, "Test Asset", resp.Name)
}

func TestCreateBrandAssetLogic_EmptyName(t *testing.T) {
	db := setupBrandTestDB(t)

	brand := createTestBrand(t, db, "Test Brand", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewCreateBrandAssetLogic(ctx, svcCtx)

	req := &types.BrandAssetReq{
		BrandId:  strconv.FormatInt(brand.Id, 10),
		Type:     "image",
		FileUrl:  "https://example.com/image.png",
		Name:     "",
		Category: "marketing",
	}

	resp, err := logic.CreateBrandAsset(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "素材名称不能为空")
}

func TestCreateBrandAssetLogic_EmptyFileUrl(t *testing.T) {
	db := setupBrandTestDB(t)

	brand := createTestBrand(t, db, "Test Brand", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewCreateBrandAssetLogic(ctx, svcCtx)

	req := &types.BrandAssetReq{
		BrandId:  strconv.FormatInt(brand.Id, 10),
		Type:     "image",
		FileUrl:  "",
		Name:     "Test Asset",
		Category: "marketing",
	}

	resp, err := logic.CreateBrandAsset(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "文件URL不能为空")
}

func TestCreateBrandAssetLogic_InvalidBrandId(t *testing.T) {
	db := setupBrandTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewCreateBrandAssetLogic(ctx, svcCtx)

	req := &types.BrandAssetReq{
		BrandId:  "invalid",
		Type:     "image",
		FileUrl:  "https://example.com/image.png",
		Name:     "Test Asset",
		Category: "marketing",
	}

	resp, err := logic.CreateBrandAsset(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "品牌ID无效")
}

func TestCreateBrandAssetLogic_BrandNotFound(t *testing.T) {
	db := setupBrandTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewCreateBrandAssetLogic(ctx, svcCtx)

	req := &types.BrandAssetReq{
		BrandId:  "999",
		Type:     "image",
		FileUrl:  "https://example.com/image.png",
		Name:     "Test Asset",
		Category: "marketing",
	}

	resp, err := logic.CreateBrandAsset(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "品牌不存在")
}

func TestUpdateBrandLogic_Success(t *testing.T) {
	db := setupBrandTestDB(t)

	brand := createTestBrand(t, db, "Test Brand", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateBrandLogic(ctx, svcCtx)

	req := &types.UpdateBrandReq{
		Id:          brand.Id,
		Name:        "Updated Brand",
		Description: "Updated description",
		Status:      "disabled",
	}

	resp, err := logic.UpdateBrand(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Updated Brand", resp.Name)
	assert.Equal(t, "disabled", resp.Status)
}

func TestUpdateBrandLogic_NotFound(t *testing.T) {
	db := setupBrandTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateBrandLogic(ctx, svcCtx)

	req := &types.UpdateBrandReq{
		Id:          999,
		Name:        "Updated Brand",
		Description: "Updated description",
	}

	resp, err := logic.UpdateBrand(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}
