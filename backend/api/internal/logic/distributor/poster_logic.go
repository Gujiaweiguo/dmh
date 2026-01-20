package distributor

import (
	"context"
	"dmh/api/internal/svc"
	"dmh/model"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

// PosterLogic 海报生成业务逻辑
type PosterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewPosterLogic 创建海报生成逻辑实例
func NewPosterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PosterLogic {
	return &PosterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GeneratePosterReq 生成海报请求
type GeneratePosterReq struct {
	Type          string `json:"type"`          // campaign/distributor
	CampaignId    *int64 `json:"campaignId"`    // 活动ID（活动海报必填）
	DistributorId *int64 `json:"distributorId"` // 分销商ID（分销商海报必填）
}

// GenerateCampaignPoster 生成活动专属海报
func (l *PosterLogic) GenerateCampaignPoster(campaignID int64, distributorID int64) (*model.PosterTemplate, error) {
	// 1. 查询活动信息
	var campaign model.Campaign
	if err := l.svcCtx.DB.Where("id = ?", campaignID).First(&campaign).Error; err != nil {
		return nil, fmt.Errorf("查询活动信息失败: %w", err)
	}

	// 2. 查询分销商信息
	var distributor model.Distributor
	if err := l.svcCtx.DB.Where("id = ?", distributorID).First(&distributor).Error; err != nil {
		return nil, fmt.Errorf("查询分销商信息失败: %w", err)
	}

	// 3. 检查是否已存在海报
	var existingPoster model.PosterTemplate
	err := l.svcCtx.DB.Where("type = ? AND campaign_id = ? AND distributor_id = ?", "campaign", campaignID, distributorID).
		First(&existingPoster).Error
	if err == nil {
		// 海报已存在，直接返回
		return &existingPoster, nil
	}
	if !errors.Is(err, errors.New("record not found")) {
		return nil, fmt.Errorf("查询海报失败: %w", err)
	}

	// 4. 生成海报URL（模拟）
	// 实际应用中，这里应该调用图片生成服务
	posterURL := fmt.Sprintf("/api/v1/posters/campaign/%d/distributor/%d", campaignID, distributorID)
	qrcodeURL := fmt.Sprintf("/api/v1/posters/qrcode?campaignId=%d&distributorId=%d", campaignID, distributorID)

	// 5. 构建海报数据
	posterData := map[string]interface{}{
		"campaignName":    campaign.Name,
		"campaignDesc":    campaign.Description,
		"distributorId":   distributorId,
		"distributorName": distributor.User.Nickname,
		"qrcodeUrl":       qrcodeURL,
		"posterUrl":       posterURL,
	}

	// 6. 创建海报记录
	poster := &model.PosterTemplate{
		Type:          "campaign",
		CampaignId:    &campaignID,
		DistributorId: &distributorID,
		TemplateUrl:   posterURL,
		PosterData:    posterData,
	}

	if err := l.svcCtx.DB.Create(poster).Error; err != nil {
		return nil, fmt.Errorf("创建海报记录失败: %w", err)
	}

	return poster, nil
}

// GenerateDistributorPoster 生成通用分销商海报
func (l *PosterLogic) GenerateDistributorPoster(distributorID int64) (*model.PosterTemplate, error) {
	// 1. 查询分销商信息
	var distributor model.Distributor
	if err := l.svcCtx.DB.Where("id = ?", distributorID).First(&distributor).Error; err != nil {
		return nil, fmt.Errorf("查询分销商信息失败: %w", err)
	}

	// 2. 查询该品牌的所有活动
	var campaigns []model.Campaign
	if err := l.svcCtx.DB.Where("brand_id = ? AND status = ?", distributor.BrandId, "active").Find(&campaigns).Error; err != nil {
		return nil, fmt.Errorf("查询活动列表失败: %w", err)
	}

	// 3. 检查是否已存在海报
	var existingPoster model.PosterTemplate
	err := l.svcCtx.DB.Where("type = ? AND distributor_id = ?", "distributor", distributorID).
		First(&existingPoster).Error
	if err == nil {
		// 海报已存在，更新数据后返回
		updatedPoster, updateErr := l.updateDistributorPoster(&existingPoster, &distributor, campaigns)
		if updateErr != nil {
			return nil, updateErr
		}
		return updatedPoster, nil
	}
	if !errors.Is(err, errors.New("record not found")) {
		return nil, fmt.Errorf("查询海报失败: %w", err)
	}

	// 4. 生成海报URL（模拟）
	posterURL := fmt.Sprintf("/api/v1/posters/distributor/%d", distributorID)
	qrcodeURL := fmt.Sprintf("/api/v1/posters/qrcode?distributorId=%d", distributorID)

	// 5. 构建海报数据
	posterData := map[string]interface{}{
		"distributorId":   distributorId,
		"distributorName": distributor.User.Nickname,
		"campaignCount":   len(campaigns),
		"campaigns":       campaigns,
		"qrcodeUrl":       qrcodeURL,
		"posterUrl":       posterURL,
	}

	// 6. 创建海报记录
	poster := &model.PosterTemplate{
		Type:          "distributor",
		DistributorId: &distributorID,
		TemplateUrl:   posterURL,
		PosterData:    posterData,
	}

	if err := l.svcCtx.DB.Create(poster).Error; err != nil {
		return nil, fmt.Errorf("创建海报记录失败: %w", err)
	}

	return poster, nil
}

// updateDistributorPoster 更新分销商海报
func (l *PosterLogic) updateDistributorPoster(poster *model.PosterTemplate, distributor *model.Distributor, campaigns []model.Campaign) (*model.PosterTemplate, error) {
	// 更新海报数据
	qrcodeURL := fmt.Sprintf("/api/v1/posters/qrcode?distributorId=%d", distributor.Id)
	posterData := map[string]interface{}{
		"distributorId":   distributor.Id,
		"distributorName": distributor.User.Nickname,
		"campaignCount":   len(campaigns),
		"campaigns":       campaigns,
		"qrcodeUrl":       qrcodeURL,
		"posterUrl":       poster.TemplateUrl,
	}

	// 更新数据库
	if err := l.svcCtx.DB.Model(poster).Updates(map[string]interface{}{
		"poster_data": posterData,
		"updated_at":  "NOW()",
	}).Error; err != nil {
		return nil, fmt.Errorf("更新海报记录失败: %w", err)
	}

	poster.PosterData = posterData
	return poster, nil
}

// GetPoster 获取海报信息
func (l *PosterLogic) GetPoster(posterID int64) (*model.PosterTemplate, error) {
	var poster model.PosterTemplate
	if err := l.svcCtx.DB.Where("id = ?", posterID).First(&poster).Error; err != nil {
		return nil, fmt.Errorf("查询海报失败: %w", err)
	}
	return &poster, nil
}
