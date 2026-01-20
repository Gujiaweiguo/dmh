package distributor

import (
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"net/http"

	"dmh/api/internal/logic/distributor"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// GetGlobalStatsHandler 获取平台全局统计数据（平台管理员）
func GetGlobalStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetGlobalStatsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewGlobalStatsLogic(r.Context(), svcCtx)
		stats, err := l.GetGlobalStats(req.StartDate, req.EndDate)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, stats)
	}
}

// GetPlatformDistributorsHandler 获取全局分销商列表（平台管理员）
func GetPlatformDistributorsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetPlatformDistributorsReq
		if err := httpx.ParseForm(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewGlobalStatsLogic(r.Context(), svcCtx)
		resp, err := l.GetPlatformDistributors(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// GetPlatformRewardsHandler 获取全局奖励列表（平台管理员）
func GetPlatformRewardsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetPlatformRewardsReq
		if err := httpx.ParseForm(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewGlobalStatsLogic(r.Context(), svcCtx)
		resp, err := l.GetPlatformRewards(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
