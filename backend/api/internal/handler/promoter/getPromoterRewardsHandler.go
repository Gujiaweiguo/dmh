package promoter

import (
	"net/http"

	"dmh/api/internal/handler/handlerutil"
	"dmh/api/internal/logic/promoter"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetPromoterRewardsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		promoterId, ok := handlerutil.ParsePathID(r, w, "/promoter/rewards/", "无效的推广员ID")
		if !ok {
			return
		}

		req := types.GetPromoterRewardsReq{
			PromoterId: promoterId,
			Page:       handlerutil.ParseQueryInt64(r, "page", 1),
			PageSize:   handlerutil.ParseQueryInt64(r, "pageSize", 20),
			Type:       handlerutil.ParseQueryString(r, "type"),
			Status:     handlerutil.ParseQueryString(r, "status"),
		}

		l := promoter.NewGetPromoterRewardsLogic(r.Context(), svcCtx)
		resp, err := l.GetPromoterRewards(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
