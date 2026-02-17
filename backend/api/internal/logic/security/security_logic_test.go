package security

import (
	"context"
	"testing"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestSecurityLogicConstructors(t *testing.T) {
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{}

	assert.NotNil(t, NewCheckPasswordStrengthLogic(ctx, svcCtx))
	assert.NotNil(t, NewForceLogoutUserLogic(ctx, svcCtx))
	assert.NotNil(t, NewGetAuditLogsLogic(ctx, svcCtx))
	assert.NotNil(t, NewGetLoginAttemptsLogic(ctx, svcCtx))
	assert.NotNil(t, NewGetPasswordPolicyLogic(ctx, svcCtx))
	assert.NotNil(t, NewGetSecurityEventsLogic(ctx, svcCtx))
	assert.NotNil(t, NewGetUserSessionsLogic(ctx, svcCtx))
	assert.NotNil(t, NewHandleSecurityEventLogic(ctx, svcCtx))
	assert.NotNil(t, NewRevokeSessionLogic(ctx, svcCtx))
	assert.NotNil(t, NewUpdatePasswordPolicyLogic(ctx, svcCtx))
}

func TestSecurityLogicCurrentMethodBehavior(t *testing.T) {
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{}

	check := NewCheckPasswordStrengthLogic(ctx, svcCtx)
	checkResp, checkErr := check.CheckPasswordStrength(&types.ChangePasswordReq{})
	assert.NoError(t, checkErr)
	assert.Nil(t, checkResp)

	force := NewForceLogoutUserLogic(ctx, svcCtx)
	forceResp, forceErr := force.ForceLogoutUser(&types.ForceLogoutReq{})
	assert.NoError(t, forceErr)
	assert.Nil(t, forceResp)

	audit := NewGetAuditLogsLogic(ctx, svcCtx)
	auditResp, auditErr := audit.GetAuditLogs()
	assert.NoError(t, auditErr)
	assert.NotNil(t, auditResp)
	assert.Equal(t, int64(0), auditResp.Total)
	assert.Len(t, auditResp.Logs, 0)

	attempts := NewGetLoginAttemptsLogic(ctx, svcCtx)
	attemptsResp, attemptsErr := attempts.GetLoginAttempts()
	assert.NoError(t, attemptsErr)
	assert.Nil(t, attemptsResp)

	policy := NewGetPasswordPolicyLogic(ctx, svcCtx)
	policyResp, policyErr := policy.GetPasswordPolicy()
	assert.NoError(t, policyErr)
	assert.Nil(t, policyResp)

	events := NewGetSecurityEventsLogic(ctx, svcCtx)
	eventsResp, eventsErr := events.GetSecurityEvents()
	assert.NoError(t, eventsErr)
	assert.Nil(t, eventsResp)

	sessions := NewGetUserSessionsLogic(ctx, svcCtx)
	sessionsResp, sessionsErr := sessions.GetUserSessions()
	assert.NoError(t, sessionsErr)
	assert.Nil(t, sessionsResp)

	handle := NewHandleSecurityEventLogic(ctx, svcCtx)
	handleResp, handleErr := handle.HandleSecurityEvent(&types.HandleSecurityEventReq{})
	assert.NoError(t, handleErr)
	assert.Nil(t, handleResp)

	revoke := NewRevokeSessionLogic(ctx, svcCtx)
	revokeResp, revokeErr := revoke.RevokeSession()
	assert.NoError(t, revokeErr)
	assert.Nil(t, revokeResp)

	update := NewUpdatePasswordPolicyLogic(ctx, svcCtx)
	updateResp, updateErr := update.UpdatePasswordPolicy(&types.UpdatePasswordPolicyReq{})
	assert.NoError(t, updateErr)
	assert.Nil(t, updateResp)
}
