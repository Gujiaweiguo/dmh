// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"context"
	"encoding/json"
	"strings"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCampaignsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCampaignsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCampaignsLogic {
	return &GetCampaignsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCampaignsLogic) GetCampaigns(req *types.GetCampaignsReq) (resp *types.CampaignListResp, err error) {
	var campaigns []model.Campaign
	var total int64

	db := l.svcCtx.DB.Model(&model.Campaign{})

	// 状态筛选
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	// 关键词搜索
	if req.Keyword != "" {
		keyword := "%" + strings.ToLower(req.Keyword) + "%"
		db = db.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", keyword, keyword)
	}

	// 获取总数
	db.Count(&total)

	// 分页
	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		db = db.Offset(int(offset)).Limit(int(req.PageSize))
	}

	// 查询数据
	if err := db.Order("id DESC").Find(&campaigns).Error; err != nil {
		return nil, err
	}

	// 获取所有品牌信息
	var brands []model.Brand
	l.svcCtx.DB.Find(&brands)
	brandMap := make(map[int64]string)
	for _, b := range brands {
		brandMap[b.Id] = b.Name
	}

	// 转换为响应格式
	campaignResps := make([]types.CampaignResp, 0, len(campaigns))
	for _, c := range campaigns {
		var formFields []types.FormField
		if c.FormFields != "" {
			json.Unmarshal([]byte(c.FormFields), &formFields)
		}

		// 解析分销奖励配置
		var distributionRewards []types.DistributorLevelRewardResp
		if c.DistributionRewards != nil {
			json.Unmarshal([]byte(*c.DistributionRewards), &distributionRewards)
		}

		campaignResps = append(campaignResps, types.CampaignResp{
			Id:                  c.Id,
			BrandId:             c.BrandId,
			BrandName:           brandMap[c.BrandId],
			Name:                c.Name,
			Description:         c.Description,
			FormFields:          formFields,
			RewardRule:          c.RewardRule,
			StartTime:           c.StartTime.Format("2006-01-02 15:04:05"),
			EndTime:             c.EndTime.Format("2006-01-02 15:04:05"),
			Status:              c.Status,
			EnableDistribution:  c.EnableDistribution,
			DistributionLevel:   c.DistributionLevel,
			DistributionRewards: distributionRewards,
			CreatedAt:           c.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	resp = &types.CampaignListResp{
		Total:     total,
		Campaigns: campaignResps,
	}

	l.Logger.Infof("获取活动列表成功，共 %d 个活动", total)

	return
}
