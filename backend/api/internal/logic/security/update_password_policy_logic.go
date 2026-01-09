package security

import (
	"context"
	"time"

	"dmh/api/internal/service"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePasswordPolicyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePasswordPolicyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePasswordPolicyLogic {
	return &UpdatePasswordPolicyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePasswordPolicyLogic) UpdatePasswordPolicy(req *types.UpdatePasswordPolicyReq) (resp *types.PasswordPolicyResp, err error) {
	// 获取当前策略
	currentPolicy, err := l.svcCtx.PasswordService.GetPasswordPolicy()
	if err != nil {
		logx.Errorf("获取当前密码策略失败: %v", err)
		return nil, err
	}

	// 更新策略字段
	policy := &model.PasswordPolicy{
		ID:                    currentPolicy.ID,
		MinLength:             currentPolicy.MinLength,
		RequireUppercase:      currentPolicy.RequireUppercase,
		RequireLowercase:      currentPolicy.RequireLowercase,
		RequireNumbers:        currentPolicy.RequireNumbers,
		RequireSpecialChars:   currentPolicy.RequireSpecialChars,
		MaxAge:                currentPolicy.MaxAge,
		HistoryCount:          currentPolicy.HistoryCount,
		MaxLoginAttempts:      currentPolicy.MaxLoginAttempts,
		LockoutDuration:       currentPolicy.LockoutDuration,
		SessionTimeout:        currentPolicy.SessionTimeout,
		MaxConcurrentSessions: currentPolicy.MaxConcurrentSessions,
		CreatedAt:             currentPolicy.CreatedAt,
		UpdatedAt:             time.Now(),
	}

	// 应用更新
	if req.MinLength != 0 {
		policy.MinLength = req.MinLength
	}
	if req.RequireUppercase {
		policy.RequireUppercase = req.RequireUppercase
	}
	if req.RequireLowercase {
		policy.RequireLowercase = req.RequireLowercase
	}
	if req.RequireNumbers {
		policy.RequireNumbers = req.RequireNumbers
	}
	if req.RequireSpecialChars {
		policy.RequireSpecialChars = req.RequireSpecialChars
	}
	if req.MaxAge != 0 {
		policy.MaxAge = req.MaxAge
	}
	if req.HistoryCount != 0 {
		policy.HistoryCount = req.HistoryCount
	}
	if req.MaxLoginAttempts != 0 {
		policy.MaxLoginAttempts = req.MaxLoginAttempts
	}
	if req.LockoutDuration != 0 {
		policy.LockoutDuration = req.LockoutDuration
	}
	if req.SessionTimeout != 0 {
		policy.SessionTimeout = req.SessionTimeout
	}
	if req.MaxConcurrentSessions != 0 {
		policy.MaxConcurrentSessions = req.MaxConcurrentSessions
	}

	// 保存更新
	err = l.svcCtx.PasswordService.UpdatePasswordPolicy(policy)
	if err != nil {
		logx.Errorf("更新密码策略失败: %v", err)
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
		"update_password_policy",
		"password_policy",
		"1",
		req,
	)

	resp = &types.PasswordPolicyResp{
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
	}

	return resp, nil
}