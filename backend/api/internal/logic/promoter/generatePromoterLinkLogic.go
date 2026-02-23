package promoter

import (
	"context"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type GeneratePromoterLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGeneratePromoterLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GeneratePromoterLinkLogic {
	return &GeneratePromoterLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GeneratePromoterLinkLogic) GeneratePromoterLink(req *types.GeneratePromoterLinkReq) (resp *types.GeneratePromoterLinkResp, err error) {
	var promoter model.Promoter
	if err := l.svcCtx.DB.First(&promoter, req.PromoterId).Error; err != nil {
		return nil, fmt.Errorf("promoter not found: %w", err)
	}

	var campaign model.Campaign
	if err := l.svcCtx.DB.First(&campaign, req.CampaignId).Error; err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	linkCode := uuid.New().String()[:8]

	link := model.PromoterLink{
		PromoterId: req.PromoterId,
		CampaignId: req.CampaignId,
		LinkCode:   linkCode,
		ClickCount: 0,
		OrderCount: 0,
		CreatedAt:  time.Now(),
	}

	if err := l.svcCtx.DB.Create(&link).Error; err != nil {
		return nil, fmt.Errorf("failed to create link: %w", err)
	}

	linkUrl := fmt.Sprintf("https://example.com/p/%s", linkCode)

	return &types.GeneratePromoterLinkResp{
		LinkCode: linkCode,
		LinkUrl:  linkUrl,
	}, nil
}
