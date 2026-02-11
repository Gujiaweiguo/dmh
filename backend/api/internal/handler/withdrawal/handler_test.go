package withdrawal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithdrawalHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetWithdrawalsHandler(nil))
	assert.NotNil(t, GetWithdrawalHandler(nil))
	assert.NotNil(t, ApplyWithdrawalHandler(nil))
	assert.NotNil(t, ApproveWithdrawalHandler(nil))
}
