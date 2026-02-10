// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"fmt"
	"net/http"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetCampaignHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCampaignReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		id := req.Id

		var campaign model.Campaign

		if err := svcCtx.DB.First(&campaign, id).Error; err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("campaign not found: %w", err))
			return
		}

		resp := &types.CampaignResp{
			Id:                 campaign.Id,
			BrandId:            campaign.BrandId,
			Name:               campaign.Name,
			Description:        campaign.Description,
			FormFields:         campaign.FormFields,
			RewardRule:         campaign.RewardRule,
			StartTime:          campaign.StartTime.Format("2006-01-02T15:04:05"),
			EndTime:            campaign.EndTime.Format("2006-01-02T15:04:05"),
			Status:             campaign.Status,
			EnableDistribution: campaign.EnableDistribution,
			DistributionLevel:  campaign.DistributionLevel,
			DistributionRewards: func() string {
				if campaign.DistributionRewards != nil {
					return *campaign.DistributionRewards
				}
				return ""
			}(),
			PaymentConfig: func() string {
				if campaign.PaymentConfig != nil {
					return *campaign.PaymentConfig
				}
				return ""
			}(),
			PosterTemplateId: campaign.PosterTemplateId,
			CreatedAt:        campaign.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt:        campaign.UpdatedAt.Format("2006-01-02T15:04:05"),
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
