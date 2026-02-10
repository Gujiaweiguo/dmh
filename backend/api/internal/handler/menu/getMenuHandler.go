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
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/menus/")
		menuIdStr := strings.Split(path, "/")[0]
		menuId, err := strconv.ParseInt(menuIdStr, 10, 64)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errors.New("Invalid menu ID in path"))
			return
		}

		l := menu.NewGetMenuLogic(r.Context(), svcCtx)
		resp, err := l.GetMenu(menuId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
