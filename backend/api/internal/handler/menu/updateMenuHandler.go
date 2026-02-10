// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"dmh/api/internal/logic/menu"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateMenuReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/api/v1/menus/")
		menuIdStr := strings.Split(path, "/")[0]

		if menuIdStr == "" {
			httpx.ErrorCtx(r.Context(), w, errors.New("menuId is required"))
			return
		}

		menuId, err := strconv.ParseInt(menuIdStr, 10, 64)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := menu.NewUpdateMenuLogic(r.Context(), svcCtx)
		resp, err := l.UpdateMenu(menuId, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
