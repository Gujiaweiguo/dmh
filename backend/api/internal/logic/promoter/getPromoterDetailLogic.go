package promoter

import (
	"context"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPromoterDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPromoterDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPromoterDetailLogic {
	return &GetPromoterDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPromoterDetailLogic) GetPromoterDetail(req *types.GetPromoterDetailReq) (resp *types.PromoterDetailResp, err error) {
	var promoter model.Promoter
	if err := l.svcCtx.DB.Preload("User").Preload("Brand").
		First(&promoter, req.Id).Error; err != nil {
		return nil, fmt.Errorf("promoter not found: %w", err)
	}

	resp = &types.PromoterDetailResp{
		Id:             promoter.Id,
		UserId:         promoter.UserId,
		BrandId:        promoter.BrandId,
		Status:         promoter.Status,
		Level:          promoter.Level,
		TotalOrders:    promoter.TotalOrders,
		TotalRewards:   promoter.TotalRewards,
		ConversionRate: promoter.ConversionRate,
		CampaignCount:  promoter.CampaignCount,
		CreatedAt:      promoter.CreatedAt.Format("2006-01-02 15:04:05"),
		Links:          make([]types.PromoterLinkResp, 0),
	}

	if promoter.User != nil {
		resp.Username = promoter.User.Username
		resp.Phone = promoter.User.Phone
	}

	if promoter.Brand != nil {
		resp.BrandName = promoter.Brand.Name
	}

	if promoter.LastActiveAt != nil {
		resp.LastActiveAt = promoter.LastActiveAt.Format("2006-01-02 15:04:05")
	}

	var links []model.PromoterLink
	l.svcCtx.DB.Where("promoter_id = ?", promoter.Id).
		Preload("Campaign").
		Order("created_at DESC").
		Limit(10).
		Find(&links)

	for _, link := range links {
		linkResp := types.PromoterLinkResp{
			Id:         link.Id,
			CampaignId: link.CampaignId,
			LinkCode:   link.LinkCode,
			ClickCount: link.ClickCount,
			OrderCount: link.OrderCount,
			CreatedAt:  link.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if link.Campaign != nil {
			linkResp.CampaignName = link.Campaign.Name
		}
		resp.Links = append(resp.Links, linkResp)
	}

	return resp, nil
}

func (l *GetPromoterDetailLogic) GetPromoterDetailById(id int64) (resp *types.PromoterDetailResp, err error) {
	return l.GetPromoterDetail(&types.GetPromoterDetailReq{Id: id})
}
