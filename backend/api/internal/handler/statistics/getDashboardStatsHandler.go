// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package statistics

import (
	"net/http"

	"dmh/api/internal/logic/statistics"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetDashboardStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDashboardStatsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := statistics.NewGetDashboardStatsLogic(r.Context(), svcCtx)
		resp, err := l.GetDashboardStats(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
