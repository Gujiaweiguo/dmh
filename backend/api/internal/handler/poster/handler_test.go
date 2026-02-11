package poster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPosterHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetPosterTemplatesHandler(nil))
	assert.NotNil(t, GenerateCampaignPosterHandler(nil))
	assert.NotNil(t, GenerateDistributorPosterHandler(nil))
	assert.NotNil(t, GetPosterRecordsHandler(nil))
}
