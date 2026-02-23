package promoter

import (
	"net/http"

	"dmh/api/internal/handler/handlerutil"
	"dmh/api/internal/logic/promoter"
	"dmh/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetPromoterDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := handlerutil.ParsePathID(r, w, "/promoter/detail/", "无效的推广员ID")
		if !ok {
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
