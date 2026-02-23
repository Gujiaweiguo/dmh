package promoter

import (
	"context"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPromoterListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPromoterListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPromoterListLogic {
	return &GetPromoterListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPromoterListLogic) GetPromoterList(req *types.GetPromoterListReq) (resp *types.PromoterListResp, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	query := l.svcCtx.DB.Model(&model.Promoter{})

	if req.BrandId > 0 {
		query = query.Where("brand_id = ?", req.BrandId)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Keyword != "" {
		query = query.Joins("LEFT JOIN users ON users.id = promoters.user_id").
			Where("users.username LIKE ? OR users.phone LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count promoters: %w", err)
	}

	var promoters []model.Promoter
	offset := int((req.Page - 1) * req.PageSize)
	if err := query.Offset(offset).Limit(int(req.PageSize)).
		Preload("User").
		Preload("Brand").
		Order("created_at DESC").
		Find(&promoters).Error; err != nil {
		return nil, fmt.Errorf("failed to get promoters: %w", err)
	}

	resp = &types.PromoterListResp{
		Total:     total,
		Promoters: make([]types.PromoterResp, 0, len(promoters)),
	}

	for _, p := range promoters {
		promoterResp := types.PromoterResp{
			Id:             p.Id,
			UserId:         p.UserId,
			BrandId:        p.BrandId,
			Status:         p.Status,
			Level:          p.Level,
			TotalOrders:    p.TotalOrders,
			TotalRewards:   p.TotalRewards,
			ConversionRate: p.ConversionRate,
			CreatedAt:      p.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if p.User != nil {
			promoterResp.Username = p.User.Username
			promoterResp.Phone = p.User.Phone
		}

		if p.Brand != nil {
			promoterResp.BrandName = p.Brand.Name
		}

		resp.Promoters = append(resp.Promoters, promoterResp)
	}

	return resp, nil
}
