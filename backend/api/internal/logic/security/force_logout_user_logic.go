package security

import (
	"context"
	"strconv"

	"dmh/api/internal/service"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ForceLogoutUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewForceLogoutUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForceLogoutUserLogic {
	return &ForceLogoutUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ForceLogoutUserLogic) ForceLogoutUser(userId string, req *types.ForceLogoutReq) (resp *types.CommonResp, err error) {
	// 转换用户ID
	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		logx.Errorf("无效的用户ID: %s", userId)
		return nil, err
	}

	// 强制用户下线
	reason := "admin_forced_logout"
	if req.Reason != "" {
		reason = req.Reason
	}

	err = l.svcCtx.SessionService.ForceLogoutUser(userIdInt, reason)
	if err != nil {
		logx.Errorf("强制用户下线失败: %v", err)
		return nil, err
	}

	// 记录审计日志
	currentUserID := l.ctx.Value("userId").(int64)
	l.svcCtx.AuditService.LogUserAction(
		&service.AuditContext{
			UserID:    &currentUserID,
			Username:  l.ctx.Value("username").(string),
			ClientIP:  l.ctx.Value("clientIP").(string),
			UserAgent: l.ctx.Value("userAgent").(string),
		},
		"force_logout_user",
		"user",
		userId,
		map[string]interface{}{
			"target_user_id": userIdInt,
			"reason":         reason,
		},
	)

	// 记录安全事件
	l.svcCtx.AuditService.LogSecurityEvent(
		"force_logout",
		"medium",
		&userIdInt,
		"",
		l.ctx.Value("clientIP").(string),
		l.ctx.Value("userAgent").(string),
		"管理员强制用户下线",
		map[string]interface{}{
			"admin_user_id": l.ctx.Value("userId").(int64),
			"admin_username": l.ctx.Value("username").(string),
			"target_user_id": userIdInt,
			"reason":         reason,
		},
	)

	resp = &types.CommonResp{
		Message: "用户已被强制下线",
	}

	return resp, nil
}