package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetPasswordPolicyHandler(nil))
	assert.NotNil(t, UpdatePasswordPolicyHandler(nil))
	assert.NotNil(t, CheckPasswordStrengthHandler(nil))
	assert.NotNil(t, GetUserSessionsHandler(nil))
	assert.NotNil(t, RevokeSessionHandler(nil))
	assert.NotNil(t, ForceLogoutUserHandler(nil))
	assert.NotNil(t, GetAuditLogsHandler(nil))
	assert.NotNil(t, GetLoginAttemptsHandler(nil))
	assert.NotNil(t, HandleSecurityEventHandler(nil))
	assert.NotNil(t, GetSecurityEventsHandler(nil))
}
