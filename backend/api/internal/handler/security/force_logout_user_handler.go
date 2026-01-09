// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package security

import (
	"net/http"

	"dmh/api/internal/logic/security"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

func ForceLogoutUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ForceLogoutReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := security.NewForceLogoutUserLogic(r.Context(), svcCtx)
		resp, err := l.ForceLogoutUser(pathvar.Vars(r)["userId"], &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
