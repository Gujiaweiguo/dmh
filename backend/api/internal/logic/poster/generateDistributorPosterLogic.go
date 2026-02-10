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

type GenerateDistributorPosterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateDistributorPosterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateDistributorPosterLogic {
	return &GenerateDistributorPosterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateDistributorPosterLogic) GenerateDistributorPoster(req *types.GeneratePosterReq, distributorId int64) (resp *types.GeneratePosterResp, err error) {
	startTime := time.Now()

	var distributor model.Distributor
	if err := l.svcCtx.DB.First(&distributor, distributorId).Error; err != nil {
		l.Errorf("Failed to query distributor: %v", err)
		return nil, fmt.Errorf("Distributor not found")
	}

	templateId := req.TemplateId
	if templateId == 0 {
		templateId = 1
	}

	var template model.PosterTemplateConfig
	if templateId > 0 {
		if err := l.svcCtx.DB.First(&template, templateId).Error; err != nil {
			template = model.PosterTemplateConfig{
				Name: "默认模板",
			}
		}
	} else {
		template = model.PosterTemplateConfig{
			Name: "默认模板",
		}
	}

	var user model.User
	if err := l.svcCtx.DB.First(&user, distributor.UserId).Error; err == nil {
		_ = user
	}

	posterService := poster.NewService("/opt/data/posters", "http://localhost:8889/api/v1")

	posterURL, err := posterService.GenerateDistributorPoster(user.Username, 0)
	if err != nil {
		l.Errorf("Failed to generate distributor poster: %v", err)
		return nil, fmt.Errorf("Failed to generate distributor poster: %w", err)
	}

	generationTime := time.Since(startTime).Milliseconds()

	posterRecord := model.PosterRecord{
		RecordType:     "distributor",
		DistributorID:  distributor.Id,
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

	l.Infof("Distributor poster generated successfully: distributorId=%d, posterUrl=%s, time=%dms", distributor.Id, posterURL, generationTime)

	resp = &types.GeneratePosterResp{
		PosterUrl:      posterURL,
		ThumbnailUrl:   "",
		FileSize:       "",
		GenerationTime: int(generationTime),
	}

	return resp, nil
}
