// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package poster

import (
	"context"
	"fmt"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/common/poster"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateCampaignPosterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateCampaignPosterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateCampaignPosterLogic {
	return &GenerateCampaignPosterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateCampaignPosterLogic) GenerateCampaignPoster(req *types.GeneratePosterReq) (resp *types.GeneratePosterResp, err error) {
	startTime := time.Now()

	var campaign model.Campaign
	if err := l.svcCtx.DB.First(&campaign, req.Id).Error; err != nil {
		l.Errorf("Failed to query campaign: %v", err)
		return nil, fmt.Errorf("Campaign not found")
	}

	templateId := req.TemplateId
	if templateId == 0 {
		templateId = campaign.PosterTemplateId
	}

	var template model.PosterTemplateConfig
	if err := l.svcCtx.DB.First(&template, templateId).Error; err != nil {
		l.Errorf("Failed to query poster template: %v", err)
		return nil, fmt.Errorf("Poster template not found")
	}

	posterService := poster.NewService("/opt/data/posters", "http://localhost:8889/api/v1")
	qrcodeData := fmt.Sprintf("campaign_id=%d", campaign.Id)

	posterURL, err := posterService.GenerateCampaignPoster(campaign.Name, campaign.Description, "", qrcodeData)
	if err != nil {
		l.Errorf("Failed to generate poster: %v", err)
		return nil, fmt.Errorf("Failed to generate poster: %w", err)
	}

	generationTime := time.Since(startTime).Milliseconds()

	posterRecord := model.PosterRecord{
		RecordType:     "campaign",
		CampaignID:     campaign.Id,
		DistributorID:  0,
		TemplateName:   template.Name,
		PosterUrl:      posterURL,
		ThumbnailUrl:   "",
		FileSize:       "",
		GenerationTime: int(generationTime),
		DownloadCount:  0,
		ShareCount:     0,
		Status:         "success",
	}

	if err := l.svcCtx.DB.Create(&posterRecord).Error; err != nil {
		l.Errorf("Failed to save poster record: %v", err)
		return nil, fmt.Errorf("Failed to save poster record: %w", err)
	}

	l.Infof("Poster generated successfully: campaignId=%d, posterUrl=%s, time=%dms", campaign.Id, posterURL, generationTime)

	resp = &types.GeneratePosterResp{
		PosterUrl:      posterURL,
		ThumbnailUrl:   "",
		FileSize:       "",
		GenerationTime: int(generationTime),
	}

	return resp, nil
}
