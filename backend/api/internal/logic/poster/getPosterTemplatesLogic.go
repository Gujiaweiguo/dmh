// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package poster

import (
	"context"
	"fmt"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPosterTemplatesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPosterTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPosterTemplatesLogic {
	return &GetPosterTemplatesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPosterTemplatesLogic) GetPosterTemplates(req *types.GetPosterTemplatesReq) (resp *types.PosterTemplateConfigListResp, err error) {
	query := l.svcCtx.DB.Model(&model.PosterTemplateConfig{})

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Keyword != "" {
		like := fmt.Sprintf("%%%s%%", req.Keyword)
		query = query.Where("name LIKE ?", like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count poster templates: %w", err)
	}

	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		query = query.Offset(int(offset)).Limit(int(req.PageSize))
	}

	var templates []model.PosterTemplateConfig
	if err := query.Order("id desc").Find(&templates).Error; err != nil {
		return nil, fmt.Errorf("failed to query poster templates: %w", err)
	}

	respTemplates := make([]types.PosterTemplateConfigResp, 0, len(templates))
	for _, template := range templates {
		respTemplates = append(respTemplates, types.PosterTemplateConfigResp{
			Id:           template.Id,
			Name:         template.Name,
			PreviewImage: template.PreviewImage,
			Config:       template.Config,
			Status:       template.Status,
			CreatedAt:    template.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt:    template.UpdatedAt.Format("2006-01-02T15:04:05"),
		})
	}

	resp = &types.PosterTemplateConfigListResp{
		Total:     total,
		Templates: respTemplates,
	}

	return resp, nil
}
