package menu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMenuHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetMenusHandler(nil))
	assert.NotNil(t, GetMenuHandler(nil))
	assert.NotNil(t, CreateMenuHandler(nil))
	assert.NotNil(t, UpdateMenuHandler(nil))
	assert.NotNil(t, DeleteMenuHandler(nil))
	assert.NotNil(t, ConfigRoleMenusHandler(nil))
	assert.NotNil(t, GetUserMenusHandler(nil))
}
