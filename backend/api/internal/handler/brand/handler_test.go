package brand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBrandHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetBrandsHandler(nil))
	assert.NotNil(t, GetBrandHandler(nil))
	assert.NotNil(t, CreateBrandHandler(nil))
	assert.NotNil(t, UpdateBrandHandler(nil))
	assert.NotNil(t, GetBrandAssetsHandler(nil))
	assert.NotNil(t, GetBrandAssetHandler(nil))
	assert.NotNil(t, CreateBrandAssetHandler(nil))
	assert.NotNil(t, UpdateBrandAssetHandler(nil))
	assert.NotNil(t, DeleteBrandAssetHandler(nil))
	assert.NotNil(t, GetBrandStatsHandler(nil))
}
