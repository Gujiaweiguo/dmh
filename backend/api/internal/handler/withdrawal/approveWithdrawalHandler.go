// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package withdrawal

import (
	"net/http"
	"strconv"

	"dmh/api/internal/logic/withdrawal"
	"dmh/api/internal/middleware"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ApproveWithdrawalHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WithdrawalApproveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		withdrawalIdStr := r.PathValue("id")
		withdrawalId, err := strconv.ParseInt(withdrawalIdStr, 10, 64)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		adminId, err := middleware.GetUserIDFromContext(r.Context())
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := withdrawal.NewApproveWithdrawalLogic(r.Context(), svcCtx)
		resp, err := l.ApproveWithdrawal(withdrawalId, &req, adminId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
