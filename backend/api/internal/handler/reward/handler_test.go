package reward

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRewardHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetRewardsHandler(nil))
	assert.NotNil(t, GetBalanceHandler(nil))
}
