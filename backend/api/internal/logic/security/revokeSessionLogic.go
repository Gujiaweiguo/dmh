// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package security

import (
	"context"
	"errors"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

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

func (l *RevokeSessionLogic) RevokeSession(sessionID string) (resp *types.CommonResp, err error) {
	if sessionID == "" {
		return nil, errors.New("会话ID无效")
	}

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		return nil, errors.New("数据库未初始化")
	}

	var session model.UserSession
	if err := l.svcCtx.DB.Where("id = ?", sessionID).First(&session).Error; err != nil {
		return nil, errors.New("会话不存在")
	}

	if session.Status == "revoked" {
		return &types.CommonResp{Message: "会话已吊销"}, nil
	}

	if err := l.svcCtx.DB.Model(&session).Updates(map[string]interface{}{
		"status":     "revoked",
		"updated_at": time.Now(),
	}).Error; err != nil {
		l.Errorf("吊销会话失败: %v", err)
		return nil, errors.New("吊销会话失败")
	}

	return &types.CommonResp{Message: "会话已吊销"}, nil
}
