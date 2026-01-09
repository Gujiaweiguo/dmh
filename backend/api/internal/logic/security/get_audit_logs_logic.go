package security

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

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
	// 从查询参数获取分页和过滤条件
	page := 1
	pageSize := 20
	filters := make(map[string]interface{})

	// 这里应该从HTTP请求中获取查询参数，暂时使用默认值
	// 在实际实现中，需要从gin.Context或其他HTTP上下文中获取参数

	logs, total, err := l.svcCtx.AuditService.GetAuditLogs(page, pageSize, filters)
	if err != nil {
		logx.Errorf("获取审计日志失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var logList []types.AuditLogResp
	for _, log := range logs {
		logResp := types.AuditLogResp{
			Id:         log.ID,
			Username:   log.Username,
			Action:     log.Action,
			Resource:   log.Resource,
			ResourceId: log.ResourceID,
			Details:    log.Details,
			ClientIp:   log.ClientIP,
			UserAgent:  log.UserAgent,
			Status:     log.Status,
			ErrorMsg:   log.ErrorMsg,
			CreatedAt:  log.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		
		if log.UserID != nil {
			logResp.UserId = *log.UserID
		}
		
		logList = append(logList, logResp)
	}

	resp = &types.AuditLogListResp{
		Total: total,
		Logs:  logList,
	}

	return resp, nil
}