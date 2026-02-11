package admin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdminHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetUsersHandler(nil))
	assert.NotNil(t, GetUserHandler(nil))
	assert.NotNil(t, CreateUserHandler(nil))
	assert.NotNil(t, UpdateUserHandler(nil))
	assert.NotNil(t, DeleteUserHandler(nil))
	assert.NotNil(t, ResetUserPasswordHandler(nil))
	assert.NotNil(t, ManageBrandAdminRelationHandler(nil))
}
