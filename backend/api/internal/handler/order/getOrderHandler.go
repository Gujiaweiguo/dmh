// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"dmh/api/internal/logic/order"
	"dmh/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/orders/")
		orderIdStr := strings.Split(path, "/")[0]
		orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errors.New("Invalid order ID in path"))
			return
		}

		l := order.NewGetOrderLogic(r.Context(), svcCtx)
		resp, err := l.GetOrder(orderId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
