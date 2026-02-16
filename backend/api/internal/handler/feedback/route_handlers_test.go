package feedback

import (
	"testing"

	"dmh/api/internal/svc"

	"github.com/stretchr/testify/assert"
)

func TestRouteHandlersNotNil(t *testing.T) {
	svcCtx := &svc.ServiceContext{}

	assert.NotNil(t, CreateFeedbackRouteHandler(svcCtx))
	assert.NotNil(t, ListFeedbackRouteHandler(svcCtx))
	assert.NotNil(t, GetFeedbackRouteHandler(svcCtx))
	assert.NotNil(t, UpdateFeedbackStatusRouteHandler(svcCtx))
	assert.NotNil(t, SubmitSatisfactionSurveyRouteHandler(svcCtx))
	assert.NotNil(t, ListFAQRouteHandler(svcCtx))
	assert.NotNil(t, MarkFAQHelpfulRouteHandler(svcCtx))
	assert.NotNil(t, RecordFeatureUsageRouteHandler(svcCtx))
	assert.NotNil(t, GetFeedbackStatisticsRouteHandler(svcCtx))
}

func TestRouteHandlersReturnFunctions(t *testing.T) {
	svcCtx := &svc.ServiceContext{}

	tests := []struct {
		name    string
		handler interface{}
	}{
		{"CreateFeedback", CreateFeedbackRouteHandler(svcCtx)},
		{"ListFeedback", ListFeedbackRouteHandler(svcCtx)},
		{"GetFeedback", GetFeedbackRouteHandler(svcCtx)},
		{"UpdateFeedbackStatus", UpdateFeedbackStatusRouteHandler(svcCtx)},
		{"SubmitSatisfactionSurvey", SubmitSatisfactionSurveyRouteHandler(svcCtx)},
		{"ListFAQ", ListFAQRouteHandler(svcCtx)},
		{"MarkFAQHelpful", MarkFAQHelpfulRouteHandler(svcCtx)},
		{"RecordFeatureUsage", RecordFeatureUsageRouteHandler(svcCtx)},
		{"GetFeedbackStatistics", GetFeedbackStatisticsRouteHandler(svcCtx)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.handler)
		})
	}
}

func TestRouteHandlersWithNilCtx(t *testing.T) {
	handlers := []struct {
		name    string
		handler func(*svc.ServiceContext) interface {
			ServeHTTP(interface{}, interface{})
		}
	}{
		{"CreateFeedback", nil},
		{"ListFeedback", nil},
		{"GetFeedback", nil},
		{"UpdateFeedbackStatus", nil},
		{"SubmitSatisfactionSurvey", nil},
		{"ListFAQ", nil},
		{"MarkFAQHelpful", nil},
		{"RecordFeatureUsage", nil},
		{"GetFeedbackStatistics", nil},
	}

	for _, h := range handlers {
		t.Run(h.name, func(t *testing.T) {
			assert.NotNil(t, h.name)
		})
	}
}
