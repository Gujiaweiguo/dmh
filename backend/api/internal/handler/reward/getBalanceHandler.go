// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package reward

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"dmh/api/internal/logic/reward"
	"dmh/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetBalanceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/rewards/balance/")
		userIdStr := strings.Split(path, "/")[0]

		if userIdStr == "" {
			httpx.ErrorCtx(r.Context(), w, errors.New("userId is required"))
			return
		}

		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := reward.NewGetBalanceLogic(r.Context(), svcCtx)
		resp, err := l.GetBalance(userId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
