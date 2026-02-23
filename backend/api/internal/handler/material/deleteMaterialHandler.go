package material

import (
	"net/http"

	"dmh/api/internal/handler/handlerutil"
	"dmh/api/internal/logic/material"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteMaterialHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := handlerutil.ParsePathID(r, w, "/material/delete/", "无效的素材ID")
		if !ok {
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
