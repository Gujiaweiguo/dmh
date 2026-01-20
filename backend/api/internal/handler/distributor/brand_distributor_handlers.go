// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package distributor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"dmh/api/internal/logic/distributor"
	apimiddleware "dmh/api/internal/middleware"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetBrandDistributorApplicationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDistributorApplicationsReq
		var pathReq types.PathBrandIDReq
		if err := httpx.ParsePath(r, &pathReq); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if err := httpx.ParseForm(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewApproveLogic(r.Context(), svcCtx)
		resp, err := l.GetPendingApplications(pathReq.BrandID, req.Page, req.PageSize, req.Status)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func GetBrandDistributorApplicationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pathReq types.PathBrandIDApplicationIDReq
		if err := httpx.ParsePath(r, &pathReq); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewApproveLogic(r.Context(), svcCtx)
		resp, err := l.GetApplicationDetail(pathReq.BrandID, pathReq.ApplicationID)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func ApproveDistributorApplicationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !requireBrandAdmin(svcCtx, w, r) {
			return
		}

		var req types.ApproveDistributorReq
		var pathReq types.PathBrandIDApplicationIDReq
		if err := httpx.ParsePath(r, &pathReq); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewApproveLogic(r.Context(), svcCtx)
		resp, err := l.ApproveApplication(pathReq.BrandID, pathReq.ApplicationID, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func GetBrandDistributorsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDistributorsReq
		var pathReq types.PathBrandIDReq
		if err := httpx.ParsePath(r, &pathReq); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if err := httpx.ParseForm(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewManagementLogic(r.Context(), svcCtx)
		resp, err := l.GetDistributors(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func GetBrandDistributorHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pathReq types.PathBrandIDDistributorIDReq
		if err := httpx.ParsePath(r, &pathReq); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewManagementLogic(r.Context(), svcCtx)
		resp, err := l.GetDistributor(pathReq.DistributorID)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func UpdateDistributorLevelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !requireBrandAdmin(svcCtx, w, r) {
			return
		}

		var req types.UpdateDistributorLevelReq
		var pathReq types.PathDistributorIDReq
		if err := httpx.ParsePath(r, &pathReq); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewManagementLogic(r.Context(), svcCtx)
		if err := l.UpdateDistributorLevel(pathReq.DistributorID, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, types.CommonResp{Message: "级别更新成功"})
		}
	}
}

func UpdateDistributorStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !requireBrandAdmin(svcCtx, w, r) {
			return
		}

		var req types.UpdateDistributorStatusReq
		var pathReq types.PathDistributorIDReq
		if err := httpx.ParsePath(r, &pathReq); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewManagementLogic(r.Context(), svcCtx)
		if err := l.UpdateDistributorStatus(pathReq.DistributorID, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, types.CommonResp{Message: "状态更新成功"})
		}
	}
}

func GetDistributorLevelRewardsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pathReq types.PathBrandIDReq
		if err := httpx.ParsePath(r, &pathReq); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewStatisticsLogic(r.Context(), svcCtx)
		resp, err := l.GetLevelRewards(pathReq.BrandID)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func SetDistributorLevelRewardsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !requireBrandAdmin(svcCtx, w, r) {
			return
		}

		var req types.SetDistributorLevelRewardsReq
		var pathReq types.PathBrandIDReq
		if err := httpx.ParsePath(r, &pathReq); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewStatisticsLogic(r.Context(), svcCtx)
		if err := l.SetLevelRewards(pathReq.BrandID, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, types.CommonResp{Message: "奖励配置保存成功"})
		}
	}
}

func GetBrandCustomersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCustomersReq
		if err := httpx.ParseForm(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewManagementLogic(r.Context(), svcCtx)
		resp, err := l.GetCustomers(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func GetBrandRewardsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetBrandRewardsReq
		if err := httpx.ParseForm(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := distributor.NewManagementLogic(r.Context(), svcCtx)
		resp, err := l.GetBrandRewards(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func requireBrandAdmin(svcCtx *svc.ServiceContext, w http.ResponseWriter, r *http.Request) bool {
	roles, err := getRolesFromRequest(r, svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "未登录或登录已过期")
		return false
	}

	for _, role := range roles {
		if role == "platform_admin" {
			writeJSONError(w, http.StatusForbidden, "平台管理员仅可查看分销数据")
			return false
		}
	}

	for _, role := range roles {
		if role == "brand_admin" {
			return true
		}
	}

	writeJSONError(w, http.StatusForbidden, "需要品牌管理员权限")
	return false
}

func getRolesFromRequest(r *http.Request, jwtSecret string) ([]string, error) {
	authHeader := r.Header.Get("Authorization")
	if strings.TrimSpace(authHeader) == "" {
		return nil, fmt.Errorf("missing authorization header")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
		return nil, fmt.Errorf("invalid authorization header")
	}

	tokenStr := strings.TrimSpace(parts[1])
	token, err := jwt.ParseWithClaims(tokenStr, &apimiddleware.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil || token == nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*apimiddleware.JWTClaims)
	if !ok || claims == nil {
		return nil, fmt.Errorf("invalid claims")
	}
	return claims.Roles, nil
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    statusCode,
		"message": message,
	})
}
