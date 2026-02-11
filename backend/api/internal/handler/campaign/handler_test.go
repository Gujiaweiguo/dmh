package campaign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCampaignHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetCampaignsHandler(nil))
	assert.NotNil(t, GetCampaignHandler(nil))
	assert.NotNil(t, CreateCampaignHandler(nil))
	assert.NotNil(t, UpdateCampaignHandler(nil))
	assert.NotNil(t, DeleteCampaignHandler(nil))
	assert.NotNil(t, SavePageConfigHandler(nil))
	assert.NotNil(t, GetPageConfigHandler(nil))
	assert.NotNil(t, GetPaymentQrcodeHandler(nil))
}
