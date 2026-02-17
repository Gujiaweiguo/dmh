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
	assert.NotNil(t, attemptsResp)
	assert.Equal(t, int64(0), attemptsResp.Total)
	assert.Len(t, attemptsResp.Attempts, 0)

	policy := NewGetPasswordPolicyLogic(ctx, svcCtx)
	policyResp, policyErr := policy.GetPasswordPolicy()
	assert.NoError(t, policyErr)
	assert.NotNil(t, policyResp)
	assert.Equal(t, 8, policyResp.MinLength)
	assert.Equal(t, 5, policyResp.MaxLoginAttempts)

	events := NewGetSecurityEventsLogic(ctx, svcCtx)
	eventsResp, eventsErr := events.GetSecurityEvents()
	assert.NoError(t, eventsErr)
	assert.NotNil(t, eventsResp)
	assert.Equal(t, int64(0), eventsResp.Total)
	assert.Len(t, eventsResp.Events, 0)

	sessions := NewGetUserSessionsLogic(ctx, svcCtx)
	sessionsResp, sessionsErr := sessions.GetUserSessions()
	assert.NoError(t, sessionsErr)
	assert.NotNil(t, sessionsResp)
	assert.Equal(t, int64(0), sessionsResp.Total)
	assert.Len(t, sessionsResp.Sessions, 0)

	handle := NewHandleSecurityEventLogic(ctx, svcCtx)
	handleResp, handleErr := handle.HandleSecurityEvent(1, &types.HandleSecurityEventReq{})
	assert.Error(t, handleErr)
	assert.Nil(t, handleResp)

	revoke := NewRevokeSessionLogic(ctx, svcCtx)
	revokeResp, revokeErr := revoke.RevokeSession("session-id")
	assert.Error(t, revokeErr)
	assert.Nil(t, revokeResp)

	update := NewUpdatePasswordPolicyLogic(ctx, svcCtx)
	updateResp, updateErr := update.UpdatePasswordPolicy(&types.UpdatePasswordPolicyReq{})
	assert.NoError(t, updateErr)
	assert.Nil(t, updateResp)
}
