package security

import (
	"context"
	"strconv"

	"dmh/api/internal/service"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandleSecurityEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHandleSecurityEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleSecurityEventLogic {
	return &HandleSecurityEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HandleSecurityEventLogic) HandleSecurityEvent(eventId string, req *types.HandleSecurityEventReq) (resp *types.CommonResp, err error) {
	// 转换事件ID
	eventIdInt, err := strconv.ParseInt(eventId, 10, 64)
	if err != nil {
		logx.Errorf("无效的事件ID: %s", eventId)
		return nil, err
	}

	// 获取当前用户ID
	handlerID := l.ctx.Value("userId").(int64)

	// 处理安全事件
	err = l.svcCtx.AuditService.HandleSecurityEvent(eventIdInt, handlerID, req.Note)
	if err != nil {
		logx.Errorf("处理安全事件失败: %v", err)
		return nil, err
	}

	// 记录审计日志
	l.svcCtx.AuditService.LogUserAction(
		&service.AuditContext{
			UserID:    &handlerID,
			Username:  l.ctx.Value("username").(string),
			ClientIP:  l.ctx.Value("clientIP").(string),
			UserAgent: l.ctx.Value("userAgent").(string),
		},
		"handle_security_event",
		"security_event",
		eventId,
		map[string]interface{}{
			"event_id": eventIdInt,
			"note":     req.Note,
		},
	)

	resp = &types.CommonResp{
		Message: "安全事件已处理",
	}

	return resp, nil
}