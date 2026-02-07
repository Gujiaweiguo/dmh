// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"fmt"
	"net/http"

	"dmh/api/internal/logic/order"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ScanOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 手动解析 query 参数
		code := r.URL.Query().Get("code")
		if code == "" {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("code parameter is required"))
			return
		}

		l := order.NewScanOrderLogic(r.Context(), svcCtx)
		resp, err := l.ScanOrder(&types.ScanOrderReq{Code: code})
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
