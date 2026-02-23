package promoter

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"dmh/api/internal/logic/promoter"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetPromoterRewardsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/promoter/rewards/")
		promoterIdStr := strings.Split(path, "/")[0]
		promoterId, err := strconv.ParseInt(promoterIdStr, 10, 64)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errors.New("无效的推广员ID"))
			return
		}

		req := types.GetPromoterRewardsReq{
			PromoterId: promoterId,
			Page:       1,
			PageSize:   20,
		}

		if pageStr := r.URL.Query().Get("page"); pageStr != "" {
			if page, err := strconv.ParseInt(pageStr, 10, 64); err == nil && page > 0 {
				req.Page = page
			}
		}
		if pageSizeStr := r.URL.Query().Get("pageSize"); pageSizeStr != "" {
			if pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64); err == nil && pageSize > 0 {
				req.PageSize = pageSize
			}
		}
		req.Type = r.URL.Query().Get("type")
		req.Status = r.URL.Query().Get("status")

		l := promoter.NewGetPromoterRewardsLogic(r.Context(), svcCtx)
		resp, err := l.GetPromoterRewards(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
