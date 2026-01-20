package distributor

import (
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"net/http"

	"dmh/api/internal/logic/distributor"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// CalculateMultiLevelRewardHandler 计算并分配多级奖励
func CalculateMultiLevelRewardHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CalculateMultiLevelRewardReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewMultiLevelRewardLogic(r.Context(), svcCtx)
		err := l.CalculateAndDistributeRewards(req.OrderId, req.CampaignId, 0, req.ReferrerId, req.Amount)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, types.CommonResp{Message: "奖励分配成功"})
	}
}
