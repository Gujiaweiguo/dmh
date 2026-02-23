package material

import (
	"context"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadMaterialLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadMaterialLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadMaterialLogic {
	return &UploadMaterialLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadMaterialLogic) UploadMaterial(req *types.UploadMaterialReq, r *http.Request) (resp *types.UploadMaterialResp, err error) {
	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("no file uploaded: %w", err)
	}
	defer file.Close()

	materialType := req.Type
	if materialType == "" {
		materialType = "image"
	}

	var url string
	var content string

	if materialType == "image" {
		url, err = l.saveFile(file, header)
		if err != nil {
			return nil, fmt.Errorf("failed to save file: %w", err)
		}
	} else {
		contentBytes, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		content = string(contentBytes)
	}

	material := model.Material{
		Name:        header.Filename,
		Description: "",
		Type:        materialType,
		Url:         url,
		Content:     content,
		CreatedAt:   time.Now(),
	}

	if err := l.svcCtx.DB.Create(&material).Error; err != nil {
		return nil, fmt.Errorf("failed to save material: %w", err)
	}

	return &types.UploadMaterialResp{
		Id:          material.Id,
		Name:        material.Name,
		Url:         material.Url,
		Type:        material.Type,
		Description: material.Description,
	}, nil
}

func (l *UploadMaterialLogic) saveFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	return fmt.Sprintf("/uploads/%s", filename), nil
}

type GetMaterialListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMaterialListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMaterialListLogic {
	return &GetMaterialListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMaterialListLogic) GetMaterialList(req *types.GetMaterialListReq) (resp *types.MaterialListResp, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	query := l.svcCtx.DB.Model(&model.Material{})

	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count materials: %w", err)
	}

	var materials []model.Material
	offset := int((req.Page - 1) * req.PageSize)
	if err := query.Offset(offset).Limit(int(req.PageSize)).
		Order("created_at DESC").
		Find(&materials).Error; err != nil {
		return nil, fmt.Errorf("failed to get materials: %w", err)
	}

	resp = &types.MaterialListResp{
		Total:     total,
		Materials: make([]types.MaterialResp, 0, len(materials)),
	}

	for _, m := range materials {
		resp.Materials = append(resp.Materials, types.MaterialResp{
			Id:          m.Id,
			Name:        m.Name,
			Description: m.Description,
			Type:        m.Type,
			Url:         m.Url,
			Content:     m.Content,
			CreatedAt:   m.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}

type DeleteMaterialLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMaterialLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMaterialLogic {
	return &DeleteMaterialLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMaterialLogic) DeleteMaterial(req *types.DeleteMaterialReq) (resp *types.DeleteMaterialResp, err error) {
	result := l.svcCtx.DB.Delete(&model.Material{}, req.Id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to delete material: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("material not found")
	}

	return &types.DeleteMaterialResp{Success: true}, nil
}
