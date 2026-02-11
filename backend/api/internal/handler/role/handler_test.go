package role

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetRolesHandler(nil))
	assert.NotNil(t, GetPermissionsHandler(nil))
	assert.NotNil(t, ConfigRolePermissionsHandler(nil))
	assert.NotNil(t, GetUserPermissionsHandler(nil))
}
