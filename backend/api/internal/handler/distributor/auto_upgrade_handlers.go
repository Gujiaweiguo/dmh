package distributor

import (
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"net/http"

	"dmh/api/internal/logic/distributor"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// AutoUpgradeDistributorHandler 自动升级为分销商
func AutoUpgradeDistributorHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AutoUpgradeDistributorReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewAutoUpgradeLogic(r.Context(), svcCtx)
		distributor, err := l.CheckAndAutoUpgradeWithCampaign(req.UserId, req.BrandId, req.OrderId, req.ReferrerId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		resp := &types.AutoUpgradeDistributorResp{
			BecomeDistributor: distributor != nil,
		}
		if distributor != nil {
			resp.DistributorId = distributor.Id
			resp.Message = "自动成为分销商成功"
		} else {
			resp.Message = "已是分销商"
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
