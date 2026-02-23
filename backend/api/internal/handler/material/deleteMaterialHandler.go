package material

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"dmh/api/internal/logic/material"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteMaterialHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/material/delete/")
		idStr := strings.Split(path, "/")[0]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errors.New("无效的素材ID"))
			return
		}

		req := types.DeleteMaterialReq{Id: id}

		l := material.NewDeleteMaterialLogic(r.Context(), svcCtx)
		resp, err := l.DeleteMaterial(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
