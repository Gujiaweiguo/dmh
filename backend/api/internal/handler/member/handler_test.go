package member

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemberHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetMembersHandler(nil))
	assert.NotNil(t, GetMemberHandler(nil))
	assert.NotNil(t, UpdateMemberHandler(nil))
	assert.NotNil(t, UpdateMemberStatusHandler(nil))
	assert.NotNil(t, GetMemberProfileHandler(nil))
}
