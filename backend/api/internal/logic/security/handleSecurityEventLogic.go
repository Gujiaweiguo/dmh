// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package security

import (
	"context"
	"errors"
	"strings"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

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

func (l *HandleSecurityEventLogic) HandleSecurityEvent(eventID int64, req *types.HandleSecurityEventReq) (resp *types.CommonResp, err error) {
	if eventID <= 0 {
		return nil, errors.New("事件ID无效")
	}

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		return nil, errors.New("数据库未初始化")
	}

	var event model.SecurityEvent
	if err := l.svcCtx.DB.Where("id = ?", eventID).First(&event).Error; err != nil {
		return nil, errors.New("安全事件不存在")
	}

	if event.Handled {
		return &types.CommonResp{Message: "事件已处理"}, nil
	}

	now := time.Now()
	updates := map[string]interface{}{
		"handled":    true,
		"handled_at": now,
	}

	if userID, ok := l.ctx.Value("userId").(int64); ok && userID > 0 {
		updates["handled_by"] = userID
	}

	if req != nil {
		note := strings.TrimSpace(req.Note)
		if note != "" {
			if strings.TrimSpace(event.Details) == "" {
				updates["details"] = note
			} else {
				updates["details"] = event.Details + "\n" + note
			}
		}
	}

	if err := l.svcCtx.DB.Model(&event).Updates(updates).Error; err != nil {
		l.Errorf("处理安全事件失败: %v", err)
		return nil, errors.New("处理安全事件失败")
	}

	return &types.CommonResp{Message: "处理成功"}, nil
}
