// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"net/http"
	"strconv"
	"strings"

	"dmh/api/internal/logic/campaign"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SavePageConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PageConfigReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		pathParts := strings.Split(r.URL.Path, "/")
		campaignIdStr := pathParts[len(pathParts)-1]
		campaignId, _ := strconv.ParseInt(campaignIdStr, 10, 64)
		req.Id = campaignId

		logic := campaign.NewSavePageConfigLogic(r.Context(), svcCtx)
		resp, err := logic.SavePageConfig(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
