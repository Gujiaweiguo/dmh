package statistics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatisticsHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetDashboardStatsHandler(nil))
}
