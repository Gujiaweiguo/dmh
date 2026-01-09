package security

import (
	"context"

	"dmh/api/internal/service"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RevokeSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRevokeSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeSessionLogic {
	return &RevokeSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RevokeSessionLogic) RevokeSession(sessionId string) (resp *types.CommonResp, err error) {
	// 撤销会话
	err = l.svcCtx.SessionService.RevokeSession(sessionId, "admin_revoked")
	if err != nil {
		logx.Errorf("撤销会话失败: %v", err)
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
		"revoke_session",
		"user_session",
		sessionId,
		map[string]interface{}{
			"session_id": sessionId,
			"reason":     "admin_revoked",
		},
	)

	resp = &types.CommonResp{
		Message: "会话已成功撤销",
	}

	return resp, nil
}