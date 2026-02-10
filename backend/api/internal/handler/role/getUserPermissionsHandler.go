// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"dmh/api/internal/logic/role"
	"dmh/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetUserPermissionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/users/")
		pathParts := strings.Split(path, "/permissions")
		if len(pathParts) < 1 {
			httpx.ErrorCtx(r.Context(), w, errors.New("userId is required"))
			return
		}

		userIdStr := strings.TrimSuffix(pathParts[0], "/")
		if userIdStr == "" {
			httpx.ErrorCtx(r.Context(), w, errors.New("userId is required"))
			return
		}

		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := role.NewGetUserPermissionsLogic(r.Context(), svcCtx)
		resp, err := l.GetUserPermissions(userId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
