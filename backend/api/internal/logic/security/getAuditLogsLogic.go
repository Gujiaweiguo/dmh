// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package security

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAuditLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAuditLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuditLogsLogic {
	return &GetAuditLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAuditLogsLogic) GetAuditLogs() (resp *types.AuditLogListResp, err error) {
	if l.svcCtx == nil || l.svcCtx.DB == nil {
		return &types.AuditLogListResp{Total: 0, Logs: []types.AuditLogResp{}}, nil
	}

	var total int64
	if err := l.svcCtx.DB.Model(&model.AuditLog{}).Count(&total).Error; err != nil {
		l.Errorf("查询审计日志总数失败: %v", err)
		return nil, err
	}

	if total == 0 {
		return &types.AuditLogListResp{Total: 0, Logs: []types.AuditLogResp{}}, nil
	}

	var auditLogs []model.AuditLog
	if err := l.svcCtx.DB.Order("created_at DESC, id DESC").Find(&auditLogs).Error; err != nil {
		l.Errorf("查询审计日志失败: %v", err)
		return nil, err
	}

	logs := make([]types.AuditLogResp, 0, len(auditLogs))
	for _, auditLog := range auditLogs {
		var userID int64
		if auditLog.UserID != nil {
			userID = *auditLog.UserID
		}

		logs = append(logs, types.AuditLogResp{
			Id:         auditLog.ID,
			UserId:     userID,
			Username:   auditLog.Username,
			Action:     auditLog.Action,
			Resource:   auditLog.Resource,
			ResourceId: auditLog.ResourceID,
			Details:    auditLog.Details,
			ClientIp:   auditLog.ClientIP,
			UserAgent:  auditLog.UserAgent,
			Status:     auditLog.Status,
			ErrorMsg:   auditLog.ErrorMsg,
			CreatedAt:  auditLog.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.AuditLogListResp{Total: total, Logs: logs}, nil
}
