// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package reward

import (
	"net/http"
	"strconv"

	"dmh/api/internal/logic/reward"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetRewardsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIdStr := r.URL.Query().Get("userId")
		orderIdStr := r.URL.Query().Get("orderId")

		req := &types.GetRewardsReq{
			UserId:  0,
			OrderId: 0,
		}

		if userIdStr != "" {
			if userId, err := strconv.ParseInt(userIdStr, 10, 64); err == nil {
				req.UserId = userId
			}
		}

		if orderIdStr != "" {
			if orderId, err := strconv.ParseInt(orderIdStr, 10, 64); err == nil {
				req.OrderId = orderId
			}
		}

		l := reward.NewGetRewardsLogic(r.Context(), svcCtx)
		resp, err := l.GetRewards(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
