package distributor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistributorHandlersConstruct(t *testing.T) {
	assert.NotNil(t, DistributorApplyHandler(nil))
	assert.NotNil(t, GetDistributorApplicationsHandler(nil))
	assert.NotNil(t, GetDistributorApplicationHandler(nil))
	assert.NotNil(t, ApproveDistributorApplicationHandler(nil))
	assert.NotNil(t, GetMyDistributorStatusHandler(nil))
	assert.NotNil(t, GetMyDistributorDashboardHandler(nil))
	assert.NotNil(t, GetDistributorSubordinatesHandler(nil))
	assert.NotNil(t, GetDistributorStatisticsHandler(nil))
	assert.NotNil(t, GetDistributorByCodeHandler(nil))
	assert.NotNil(t, GenerateDistributorLinkHandler(nil))
	assert.NotNil(t, TrackDistributorLinkHandler(nil))
	assert.NotNil(t, GetDistributorLinksHandler(nil))
	assert.NotNil(t, GetDistributorQrcodeHandler(nil))
	assert.NotNil(t, GetDistributorRewardsHandler(nil))
	assert.NotNil(t, GetBrandDistributorsHandler(nil))
	assert.NotNil(t, GetBrandDistributorHandler(nil))
	assert.NotNil(t, UpdateDistributorStatusHandler(nil))
	assert.NotNil(t, GetBrandDistributorApplicationsHandler(nil))
	assert.NotNil(t, GetBrandDistributorApplicationHandler(nil))
	assert.NotNil(t, UpdateDistributorLevelHandler(nil))
	assert.NotNil(t, GetDistributorLevelRewardsHandler(nil))
	assert.NotNil(t, SetDistributorLevelRewardsHandler(nil))
}
