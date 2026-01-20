package distributor

import (
	"context"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"net/http"

	"dmh/api/internal/logic/distributor"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// ApplyWithdrawal 申请提现
func ApplyWithdrawalHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApplyWithdrawalReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewWithdrawalLogic(r.Context(), svcCtx)
		resp, err := l.ApplyWithdrawal(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// GetMyWithdrawals 获取我的提现记录
func GetMyWithdrawalsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PaginationReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewWithdrawalLogic(r.Context(), svcCtx)
		resp, err := l.GetMyWithdrawals(req.Page, req.PageSize)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// ApproveWithdrawal 审批通过提现（平台管理员）
func ApproveWithdrawalHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApproveWithdrawalReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewWithdrawalLogic(r.Context(), svcCtx)
		var pathReq types.PathIDReq
		if err := httpx.ParsePath(r, &pathReq); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		resp, err := l.ApproveWithdrawal(pathReq.ID, req.Notes)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, types.CommonResp{Message: "提现已批准"})
	}
}

// RejectWithdrawal 审批拒绝提现（平台管理员）
func RejectWithdrawalHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RejectWithdrawalReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewWithdrawalLogic(r.Context(), svcCtx)
		var pathReq types.PathIDReq
		if err := httpx.ParsePath(r, &pathReq); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		resp, err := l.RejectWithdrawal(pathReq.ID, req.Reason)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, types.CommonResp{Message: "提现已拒绝"})
	}
}

// GetPlatformWithdrawals 获取平台提现记录（平台管理员）
func GetPlatformWithdrawalsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetWithdrawalsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewWithdrawalLogic(r.Context(), svcCtx)
		resp, err := l.GetWithdrawals(req.BrandId, req.Status, req.Page, req.PageSize)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
