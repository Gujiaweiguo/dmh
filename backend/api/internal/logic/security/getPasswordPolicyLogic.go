// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package security

import (
	"context"
	"errors"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetPasswordPolicyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPasswordPolicyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPasswordPolicyLogic {
	return &GetPasswordPolicyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPasswordPolicyLogic) GetPasswordPolicy() (resp *types.PasswordPolicyResp, err error) {
	defaultResp := &types.PasswordPolicyResp{
		MinLength:             8,
		RequireUppercase:      true,
		RequireLowercase:      true,
		RequireNumbers:        true,
		RequireSpecialChars:   true,
		MaxAge:                90,
		HistoryCount:          5,
		MaxLoginAttempts:      5,
		LockoutDuration:       30,
		SessionTimeout:        480,
		MaxConcurrentSessions: 3,
	}

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		return defaultResp, nil
	}

	var policy model.PasswordPolicy
	err = l.svcCtx.DB.Order("id ASC").First(&policy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return defaultResp, nil
		}
		l.Errorf("查询密码策略失败: %v", err)
		return nil, err
	}

	return &types.PasswordPolicyResp{
		Id:                    policy.ID,
		MinLength:             policy.MinLength,
		RequireUppercase:      policy.RequireUppercase,
		RequireLowercase:      policy.RequireLowercase,
		RequireNumbers:        policy.RequireNumbers,
		RequireSpecialChars:   policy.RequireSpecialChars,
		MaxAge:                policy.MaxAge,
		HistoryCount:          policy.HistoryCount,
		MaxLoginAttempts:      policy.MaxLoginAttempts,
		LockoutDuration:       policy.LockoutDuration,
		SessionTimeout:        policy.SessionTimeout,
		MaxConcurrentSessions: policy.MaxConcurrentSessions,
		CreatedAt:             policy.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:             policy.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
