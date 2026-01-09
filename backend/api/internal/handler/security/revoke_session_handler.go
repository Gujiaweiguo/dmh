// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package security

import (
	"net/http"

	"dmh/api/internal/logic/security"
	"dmh/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

func RevokeSessionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := security.NewRevokeSessionLogic(r.Context(), svcCtx)
		resp, err := l.RevokeSession(pathvar.Vars(r)["sessionId"])
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
