package feedback

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type FeedbackHandlerIntegrationTestSuite struct {
	suite.Suite
	db *gorm.DB

	createFeedbackHandler           *CreateFeedbackHandler
	listFeedbackHandler             *ListFeedbackHandler
	getFeedbackHandler              *GetFeedbackHandler
	updateFeedbackStatusHandler     *UpdateFeedbackStatusHandler
	submitSatisfactionSurveyHandler *SubmitSatisfactionSurveyHandler
	listFAQHandler                  *ListFAQHandler
	markFAQHelpfulHandler           *MarkFAQHelpfulHandler
	recordFeatureUsageHandler       *RecordFeatureUsageHandler
	getFeedbackStatisticsHandler    *GetFeedbackStatisticsHandler
}

func (suite *FeedbackHandlerIntegrationTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	err = db.AutoMigrate(
		&model.User{},
		&model.UserFeedback{},
		&model.FeatureSatisfactionSurvey{},
		&model.FAQItem{},
		&model.FeatureUsageStat{},
		&model.FeedbackTag{},
		&model.FeedbackTagRelation{},
	)
	suite.Require().NoError(err)

	suite.db = db
	svcCtx := &svc.ServiceContext{DB: db}

	suite.createFeedbackHandler = NewCreateFeedbackHandler(svcCtx)
	suite.listFeedbackHandler = NewListFeedbackHandler(svcCtx)
	suite.getFeedbackHandler = NewGetFeedbackHandler(svcCtx)
	suite.updateFeedbackStatusHandler = NewUpdateFeedbackStatusHandler(svcCtx)
	suite.submitSatisfactionSurveyHandler = NewSubmitSatisfactionSurveyHandler(svcCtx)
	suite.listFAQHandler = NewListFAQHandler(svcCtx)
	suite.markFAQHelpfulHandler = NewMarkFAQHelpfulHandler(svcCtx)
	suite.recordFeatureUsageHandler = NewRecordFeatureUsageHandler(svcCtx)
	suite.getFeedbackStatisticsHandler = NewGetFeedbackStatisticsHandler(svcCtx)
}

func (suite *FeedbackHandlerIntegrationTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	_ = sqlDB.Close()
}

func (suite *FeedbackHandlerIntegrationTestSuite) SetupTest() {
	suite.db.Exec("DELETE FROM user_feedbacks")
	suite.db.Exec("DELETE FROM feature_satisfaction_surveys")
	suite.db.Exec("DELETE FROM faq_items")
	suite.db.Exec("DELETE FROM feature_usage_stats")
	suite.db.Exec("DELETE FROM feedback_tags")
	suite.db.Exec("DELETE FROM feedback_tag_relations")
	suite.db.Exec("DELETE FROM users")
}

func (suite *FeedbackHandlerIntegrationTestSuite) createTestUser(id int64, username string) {
	user := &model.User{
		Id:       id,
		Username: username,
		Password: "hashedpassword",
		Phone:    fmt.Sprintf("138%08d", id),
		Status:   "active",
	}
	_ = suite.db.Create(user).Error
}

func withAuth(req *http.Request, userID int64, role string) *http.Request {
	ctx := context.WithValue(req.Context(), "userId", userID)
	ctx = context.WithValue(ctx, "userRole", role)
	return req.WithContext(ctx)
}

func (suite *FeedbackHandlerIntegrationTestSuite) TestCreateListGetFeedbackFlow() {
	suite.createTestUser(1, "u1")
	suite.createTestUser(2, "u2")

	createReq := types.CreateFeedbackReq{
		Category: "poster",
		Title:    "反馈标题",
		Content:  "反馈内容",
	}
	body, _ := json.Marshal(createReq)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/feedback", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	suite.createFeedbackHandler.CreateFeedback(rec, withAuth(req, 1, "participant"))
	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var created types.FeedbackResp
	err := json.Unmarshal(rec.Body.Bytes(), &created)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), created.UserId)
	assert.Equal(suite.T(), "pending", created.Status)

	_ = suite.db.Create(&model.UserFeedback{UserID: 2, Category: "payment", Title: "other", Content: "other", Status: "pending"}).Error

	listRec := httptest.NewRecorder()
	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/feedback?page=1&pageSize=10", nil)
	suite.listFeedbackHandler.ListFeedback(listRec, withAuth(listReq, 1, "participant"))
	assert.Equal(suite.T(), http.StatusOK, listRec.Code)

	var listResp types.FeedbackListResp
	err = json.Unmarshal(listRec.Body.Bytes(), &listResp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), listResp.Total)

	adminListRec := httptest.NewRecorder()
	adminListReq := httptest.NewRequest(http.MethodGet, "/api/v1/feedback?page=1&pageSize=10", nil)
	suite.listFeedbackHandler.ListFeedback(adminListRec, withAuth(adminListReq, 1, "platform_admin"))
	assert.Equal(suite.T(), http.StatusOK, adminListRec.Code)

	getRec := httptest.NewRecorder()
	getReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/feedback?id=%d", created.Id), nil)
	suite.getFeedbackHandler.GetFeedback(getRec, withAuth(getReq, 1, "participant"))
	assert.Equal(suite.T(), http.StatusOK, getRec.Code)

	forbiddenRec := httptest.NewRecorder()
	forbiddenReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/feedback?id=%d", created.Id), nil)
	suite.getFeedbackHandler.GetFeedback(forbiddenRec, withAuth(forbiddenReq, 2, "participant"))
	assert.NotEqual(suite.T(), http.StatusOK, forbiddenRec.Code)
}

func (suite *FeedbackHandlerIntegrationTestSuite) TestUpdateStatusAndSurveyFlow() {
	suite.createTestUser(1, "u1")
	suite.createTestUser(2, "admin")

	feedback := &model.UserFeedback{UserID: 1, Category: "poster", Title: "待处理", Content: "content", Status: "pending"}
	_ = suite.db.Create(feedback).Error

	forbiddenReqBody := []byte(fmt.Sprintf(`{"id":%d,"status":"reviewing"}`, feedback.ID))
	forbiddenRec := httptest.NewRecorder()
	forbiddenReq := httptest.NewRequest(http.MethodPost, "/api/v1/feedback/status", bytes.NewBuffer(forbiddenReqBody))
	forbiddenReq.Header.Set("Content-Type", "application/json")
	suite.updateFeedbackStatusHandler.UpdateFeedbackStatus(forbiddenRec, withAuth(forbiddenReq, 1, "participant"))
	assert.NotEqual(suite.T(), http.StatusOK, forbiddenRec.Code)

	updateReqBody := []byte(fmt.Sprintf(`{"id":%d,"status":"resolved","response":"done"}`, feedback.ID))
	updateRec := httptest.NewRecorder()
	updateReq := httptest.NewRequest(http.MethodPost, "/api/v1/feedback/status", bytes.NewBuffer(updateReqBody))
	updateReq.Header.Set("Content-Type", "application/json")
	suite.updateFeedbackStatusHandler.UpdateFeedbackStatus(updateRec, withAuth(updateReq, 2, "platform_admin"))
	assert.Equal(suite.T(), http.StatusOK, updateRec.Code)

	var refreshed model.UserFeedback
	err := suite.db.First(&refreshed, feedback.ID).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "resolved", refreshed.Status)
	assert.NotNil(suite.T(), refreshed.ResolvedAt)

	surveyReq := []byte(`{"feature":"poster","easeOfUse":5,"performance":4,"reliability":5,"overallSatisfaction":5,"wouldRecommend":5}`)
	surveyRec := httptest.NewRecorder()
	surveyHttpReq := httptest.NewRequest(http.MethodPost, "/api/v1/feedback/survey", bytes.NewBuffer(surveyReq))
	surveyHttpReq.Header.Set("Content-Type", "application/json")
	suite.submitSatisfactionSurveyHandler.SubmitSatisfactionSurvey(surveyRec, withAuth(surveyHttpReq, 1, "participant"))
	assert.Equal(suite.T(), http.StatusOK, surveyRec.Code)
}

func (suite *FeedbackHandlerIntegrationTestSuite) TestFAQUsageAndStatisticsFlow() {
	suite.createTestUser(1, "u1")

	faq := &model.FAQItem{Category: "poster", Question: "Q", Answer: "A", SortOrder: 1, HelpfulCount: 0}
	_ = suite.db.Create(faq).Error

	listFAQRec := httptest.NewRecorder()
	listFAQReq := httptest.NewRequest(http.MethodGet, "/api/v1/faq?category=poster", nil)
	suite.listFAQHandler.ListFAQ(listFAQRec, listFAQReq)
	assert.Equal(suite.T(), http.StatusOK, listFAQRec.Code)

	markReq := []byte(fmt.Sprintf(`{"id":%d,"type":"helpful"}`, faq.ID))
	markRec := httptest.NewRecorder()
	markHTTPReq := httptest.NewRequest(http.MethodPost, "/api/v1/faq/helpful", bytes.NewBuffer(markReq))
	markHTTPReq.Header.Set("Content-Type", "application/json")
	suite.markFAQHelpfulHandler.MarkFAQHelpful(markRec, markHTTPReq)
	assert.Equal(suite.T(), http.StatusOK, markRec.Code)

	usageReq := []byte(`{"feature":"poster","action":"generate","success":true,"durationMs":1200}`)
	usageRec := httptest.NewRecorder()
	usageHTTPReq := httptest.NewRequest(http.MethodPost, "/api/v1/feedback/usage", bytes.NewBuffer(usageReq))
	usageHTTPReq.Header.Set("Content-Type", "application/json")
	suite.recordFeatureUsageHandler.RecordFeatureUsage(usageRec, withAuth(usageHTTPReq, 1, "participant"))
	assert.Equal(suite.T(), http.StatusOK, usageRec.Code)

	_ = suite.db.Create(&model.UserFeedback{UserID: 1, Category: "poster", Title: "f1", Content: "c1", Status: "pending", Rating: intPtr(4), CreatedAt: time.Now()}).Error
	_ = suite.db.Create(&model.UserFeedback{UserID: 1, Category: "payment", Title: "f2", Content: "c2", Status: "resolved", Rating: intPtr(5), CreatedAt: time.Now()}).Error

	forbiddenStatsRec := httptest.NewRecorder()
	forbiddenStatsReq := httptest.NewRequest(http.MethodGet, "/api/v1/feedback/statistics", nil)
	suite.getFeedbackStatisticsHandler.GetFeedbackStatistics(forbiddenStatsRec, withAuth(forbiddenStatsReq, 1, "participant"))
	assert.NotEqual(suite.T(), http.StatusOK, forbiddenStatsRec.Code)

	statsRec := httptest.NewRecorder()
	statsReq := httptest.NewRequest(http.MethodGet, "/api/v1/feedback/statistics", nil)
	suite.getFeedbackStatisticsHandler.GetFeedbackStatistics(statsRec, withAuth(statsReq, 1, "platform_admin"))
	assert.Equal(suite.T(), http.StatusOK, statsRec.Code)

	var statsResp types.FeedbackStatisticsResp
	err := json.Unmarshal(statsRec.Body.Bytes(), &statsResp)
	assert.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), statsResp.TotalFeedbacks, int64(2))
	assert.Greater(suite.T(), statsResp.AverageRating, 0.0)
}

func TestFeedbackHandlerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(FeedbackHandlerIntegrationTestSuite))
}

func intPtr(v int) *int {
	return &v
}
