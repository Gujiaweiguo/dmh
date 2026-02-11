package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetSyncHealthHandler(nil))
	assert.NotNil(t, GetSyncStatusHandler(nil))
	assert.NotNil(t, GetSyncStatsHandler(nil))
	assert.NotNil(t, RetrySynHandler(nil))
}
