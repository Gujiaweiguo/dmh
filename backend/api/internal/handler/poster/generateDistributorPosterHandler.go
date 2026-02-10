// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package poster

import (
	"errors"
	"net/http"
	"strconv"

	"dmh/api/internal/logic/poster"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GenerateDistributorPosterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GeneratePosterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		path := r.URL.Path
		prefix := "/api/v1/distributors/"
		suffix := "/poster"
		if len(path) > len(prefix)+len(suffix) {
			idStr := path[len(prefix) : len(path)-len(suffix)]
			distributorId, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}

			l := poster.NewGenerateDistributorPosterLogic(r.Context(), svcCtx)
			resp, err := l.GenerateDistributorPoster(&req, distributorId)
			if err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
			} else {
				httpx.OkJsonCtx(r.Context(), w, resp)
			}
		} else {
			httpx.ErrorCtx(r.Context(), w, errors.New("Invalid distributor ID in path"))
		}
	}
}
