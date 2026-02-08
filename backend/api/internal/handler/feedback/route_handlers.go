package feedback

import (
	"net/http"

	"dmh/api/internal/svc"
)

func CreateFeedbackRouteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	h := NewCreateFeedbackHandler(svcCtx)
	return h.CreateFeedback
}

func ListFeedbackRouteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	h := NewListFeedbackHandler(svcCtx)
	return h.ListFeedback
}

func GetFeedbackRouteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	h := NewGetFeedbackHandler(svcCtx)
	return h.GetFeedback
}

func UpdateFeedbackStatusRouteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	h := NewUpdateFeedbackStatusHandler(svcCtx)
	return h.UpdateFeedbackStatus
}

func SubmitSatisfactionSurveyRouteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	h := NewSubmitSatisfactionSurveyHandler(svcCtx)
	return h.SubmitSatisfactionSurvey
}

func ListFAQRouteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	h := NewListFAQHandler(svcCtx)
	return h.ListFAQ
}

func MarkFAQHelpfulRouteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	h := NewMarkFAQHelpfulHandler(svcCtx)
	return h.MarkFAQHelpful
}

func RecordFeatureUsageRouteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	h := NewRecordFeatureUsageHandler(svcCtx)
	return h.RecordFeatureUsage
}

func GetFeedbackStatisticsRouteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	h := NewGetFeedbackStatisticsHandler(svcCtx)
	return h.GetFeedbackStatistics
}
