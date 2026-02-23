package promoter

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"dmh/api/internal/logic/promoter"
	"dmh/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetPromoterDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path := strings.TrimPrefix(r.URL.Path, "/promoter/detail/")
		idStr := strings.Split(path, "/")[0]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errors.New("无效的推广员ID"))
			return
		}

		l := promoter.NewGetPromoterDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetPromoterDetailById(id)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
