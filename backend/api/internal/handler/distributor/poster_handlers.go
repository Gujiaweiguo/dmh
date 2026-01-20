package distributor

import (
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"net/http"

	"dmh/api/internal/logic/distributor"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// GenerateCampaignPosterHandler 生成活动专属海报
func GenerateCampaignPosterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GeneratePosterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewPosterLogic(r.Context(), svcCtx)
		poster, err := l.GenerateCampaignPoster(req.CampaignId, getUserId(r))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, types.GeneratePosterResp{
			PosterUrl:  poster.TemplateUrl,
			LinkCode:   "",
			IsGeneric:  false,
			CampaignId: req.CampaignId,
		})
	}
}

// GenerateDistributorPosterHandler 生成通用分销商海报
func GenerateDistributorPosterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := distributor.NewPosterLogic(r.Context(), svcCtx)
		poster, err := l.GenerateDistributorPoster(getUserId(r))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, types.GeneratePosterResp{
			PosterUrl:  poster.TemplateUrl,
			LinkCode:   "",
			IsGeneric:  true,
			CampaignId: 0,
		})
	}
}

// getUserId 从请求上下文中获取用户ID
func getUserId(r *http.Request) int64 {
	return r.Context().Value("userId").(int64)
}
