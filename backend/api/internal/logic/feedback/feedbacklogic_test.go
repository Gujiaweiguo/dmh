package feedback

import (
	"context"
	"fmt"
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

type FeedbackLogicTestSuite struct {
	suite.Suite
	db     *gorm.DB
	svcCtx *svc.ServiceContext
}

func (suite *FeedbackLogicTestSuite) SetupSuite() {
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
	suite.svcCtx = &svc.ServiceContext{DB: db}
}

func (suite *FeedbackLogicTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func (suite *FeedbackLogicTestSuite) SetupTest() {
	suite.db.Exec("DELETE FROM user_feedbacks")
	suite.db.Exec("DELETE FROM feature_satisfaction_surveys")
	suite.db.Exec("DELETE FROM faq_items")
	suite.db.Exec("DELETE FROM feature_usage_stats")
	suite.db.Exec("DELETE FROM feedback_tags")
	suite.db.Exec("DELETE FROM feedback_tag_relations")
	suite.db.Exec("DELETE FROM users")
}

func (suite *FeedbackLogicTestSuite) createTestUser(id int64, username string) *model.User {
	user := &model.User{
		Id:       id,
		Username: username,
		Password: "hashedpassword",
		Phone:    fmt.Sprintf("138%08d", id),
		Status:   "active",
	}
	suite.db.Create(user)
	return user
}

func (suite *FeedbackLogicTestSuite) TestCreateFeedback_Success() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser")

	l := NewCreateFeedbackLogic(ctx, suite.svcCtx)
	rating := 5
	req := &types.CreateFeedbackReq{
		Category:    "poster",
		Rating:      &rating,
		Title:       "测试反馈",
		Content:     "这是一个测试反馈内容",
		Priority:    "high",
		DeviceInfo:  "iPhone 14",
		BrowserInfo: "Safari",
	}

	resp, err := l.CreateFeedback(req, 1)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(1), resp.UserId)
	assert.Equal(suite.T(), "poster", resp.Category)
	assert.Equal(suite.T(), "测试反馈", resp.Title)
	assert.Equal(suite.T(), "high", resp.Priority)
	assert.Equal(suite.T(), "pending", resp.Status)
}

func (suite *FeedbackLogicTestSuite) TestCreateFeedback_Validation_EmptyTitle() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser")

	l := NewCreateFeedbackLogic(ctx, suite.svcCtx)
	req := &types.CreateFeedbackReq{
		Category: "poster",
		Title:    "",
		Content:  "test content",
	}

	resp, err := l.CreateFeedback(req, 1)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "标题")
}

func (suite *FeedbackLogicTestSuite) TestCreateFeedback_Validation_EmptyContent() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser")

	l := NewCreateFeedbackLogic(ctx, suite.svcCtx)
	req := &types.CreateFeedbackReq{
		Category: "poster",
		Title:    "test title",
		Content:  "",
	}

	resp, err := l.CreateFeedback(req, 1)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "内容")
}

func (suite *FeedbackLogicTestSuite) TestCreateFeedback_Validation_EmptyCategory() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser")

	l := NewCreateFeedbackLogic(ctx, suite.svcCtx)
	req := &types.CreateFeedbackReq{
		Title:   "test title",
		Content: "test content",
	}

	resp, err := l.CreateFeedback(req, 1)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "类别")
}

func (suite *FeedbackLogicTestSuite) TestCreateFeedback_Validation_InvalidRating() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser")

	l := NewCreateFeedbackLogic(ctx, suite.svcCtx)
	rating := 10
	req := &types.CreateFeedbackReq{
		Category: "poster",
		Rating:   &rating,
		Title:    "test title",
		Content:  "test content",
	}

	resp, err := l.CreateFeedback(req, 1)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "评分")
}

func (suite *FeedbackLogicTestSuite) TestCreateFeedback_DefaultPriority() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser")

	l := NewCreateFeedbackLogic(ctx, suite.svcCtx)
	req := &types.CreateFeedbackReq{
		Category: "poster",
		Title:    "test title",
		Content:  "test content",
	}

	resp, err := l.CreateFeedback(req, 1)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "medium", resp.Priority)
}

func (suite *FeedbackLogicTestSuite) TestListFeedback_UserViewOwnFeedbacks() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser")

	feedback1 := &model.UserFeedback{UserID: 1, Category: "poster", Title: "Feedback 1", Content: "Content 1", Status: "pending"}
	feedback2 := &model.UserFeedback{UserID: 2, Category: "payment", Title: "Feedback 2", Content: "Content 2", Status: "pending"}
	suite.db.Create(feedback1)
	suite.db.Create(feedback2)

	l := NewListFeedbackLogic(ctx, suite.svcCtx)
	req := &types.ListFeedbackReq{}

	resp, err := l.ListFeedback(req, 1, "participant")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(1), resp.Total)
	assert.Len(suite.T(), resp.Feedbacks, 1)
}

func (suite *FeedbackLogicTestSuite) TestListFeedback_AdminViewAllFeedbacks() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser")

	for i := 1; i <= 3; i++ {
		feedback := &model.UserFeedback{UserID: int64(i), Category: "poster", Title: "Feedback", Content: "Content", Status: "pending"}
		suite.db.Create(feedback)
	}

	l := NewListFeedbackLogic(ctx, suite.svcCtx)
	req := &types.ListFeedbackReq{}

	resp, err := l.ListFeedback(req, 1, "platform_admin")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(3), resp.Total)
	assert.Len(suite.T(), resp.Feedbacks, 3)
}

func (suite *FeedbackLogicTestSuite) TestListFeedback_FilterByStatus() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser")

	suite.db.Create(&model.UserFeedback{UserID: 1, Category: "poster", Title: "F1", Content: "C1", Status: "pending"})
	suite.db.Create(&model.UserFeedback{UserID: 1, Category: "poster", Title: "F2", Content: "C2", Status: "resolved"})
	suite.db.Create(&model.UserFeedback{UserID: 1, Category: "poster", Title: "F3", Content: "C3", Status: "pending"})

	l := NewListFeedbackLogic(ctx, suite.svcCtx)
	req := &types.ListFeedbackReq{Status: "pending"}

	resp, err := l.ListFeedback(req, 1, "platform_admin")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(2), resp.Total)
}

func (suite *FeedbackLogicTestSuite) TestListFeedback_Pagination() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser")

	for i := 1; i <= 5; i++ {
		feedback := &model.UserFeedback{UserID: 1, Category: "poster", Title: "Feedback", Content: "Content", Status: "pending"}
		suite.db.Create(feedback)
	}

	l := NewListFeedbackLogic(ctx, suite.svcCtx)
	req := &types.ListFeedbackReq{Page: 1, PageSize: 2}

	resp, err := l.ListFeedback(req, 1, "platform_admin")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(5), resp.Total)
	assert.Len(suite.T(), resp.Feedbacks, 2)
}

func (suite *FeedbackLogicTestSuite) TestGetFeedback_Success() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser")

	feedback := &model.UserFeedback{UserID: 1, Category: "poster", Title: "Test Feedback", Content: "Test Content", Status: "pending"}
	suite.db.Create(feedback)

	l := NewGetFeedbackLogic(ctx, suite.svcCtx)
	resp, err := l.GetFeedback(feedback.ID, 1, "participant")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), feedback.ID, resp.Id)
	assert.Equal(suite.T(), "Test Feedback", resp.Title)
}

func (suite *FeedbackLogicTestSuite) TestGetFeedback_AccessDenied() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser")
	suite.createTestUser(2, "testuser2")

	feedback := &model.UserFeedback{UserID: 2, Category: "poster", Title: "Test Feedback", Content: "Test Content", Status: "pending"}
	suite.db.Create(feedback)

	l := NewGetFeedbackLogic(ctx, suite.svcCtx)
	resp, err := l.GetFeedback(feedback.ID, 1, "participant")

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "无权")
}

func (suite *FeedbackLogicTestSuite) TestUpdateFeedbackStatus_Success() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser")
	suite.createTestUser(2, "admin")

	assigneeId := int64(2)
	feedback := &model.UserFeedback{UserID: 1, Category: "poster", Title: "Test Feedback", Content: "Test Content", Status: "pending", AssigneeID: nil}
	suite.db.Create(feedback)

	l := NewUpdateFeedbackStatusLogic(ctx, suite.svcCtx)
	req := &types.UpdateFeedbackStatusReq{Id: feedback.ID, Status: "reviewing", AssigneeId: &assigneeId, Response: "我们正在处理"}

	resp, err := l.UpdateFeedbackStatus(req, 2, "platform_admin")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "reviewing", resp.Status)
	assert.Equal(suite.T(), "我们正在处理", resp.Response)
	assert.Equal(suite.T(), &assigneeId, resp.AssigneeId)
}

func (suite *FeedbackLogicTestSuite) TestUpdateFeedbackStatus_ResolvedTime() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser")
	suite.createTestUser(2, "admin")

	assigneeId := int64(2)
	feedback := &model.UserFeedback{UserID: 1, Category: "poster", Title: "Test Feedback", Content: "Test Content", Status: "pending", AssigneeID: nil}
	suite.db.Create(feedback)

	l := NewUpdateFeedbackStatusLogic(ctx, suite.svcCtx)
	req := &types.UpdateFeedbackStatusReq{Id: feedback.ID, Status: "resolved", AssigneeId: &assigneeId, Response: "问题已解决"}

	resp, err := l.UpdateFeedbackStatus(req, 2, "platform_admin")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "resolved", resp.Status)
	assert.NotNil(suite.T(), resp.ResolvedAt)
}

func (suite *FeedbackLogicTestSuite) TestUpdateFeedbackStatus_PermissionDenied() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser")

	feedback := &model.UserFeedback{UserID: 1, Category: "poster", Title: "Test Feedback", Content: "Test Content", Status: "pending"}
	suite.db.Create(feedback)

	l := NewUpdateFeedbackStatusLogic(ctx, suite.svcCtx)
	req := &types.UpdateFeedbackStatusReq{Id: feedback.ID, Status: "reviewing"}

	resp, err := l.UpdateFeedbackStatus(req, 1, "participant")

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "无权")
}

func (suite *FeedbackLogicTestSuite) TestSubmitSatisfactionSurvey_Success() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser")

	easeOfUse := 5
	performance := 4
	reliability := 5
	overallSatisfaction := 5
	wouldRecommend := 5
	l := NewSubmitSatisfactionSurveyLogic(ctx, suite.svcCtx)
	req := &types.SubmitSatisfactionSurveyReq{
		Feature:                "poster",
		EaseOfUse:              &easeOfUse,
		Performance:            &performance,
		Reliability:            &reliability,
		OverallSatisfaction:    &overallSatisfaction,
		WouldRecommend:         &wouldRecommend,
		MostLiked:              "拖拽功能",
		LeastLiked:             "文本编辑",
		ImprovementSuggestions: "增加更多模板",
		WouldLikeMoreFeatures:  "视频支持",
	}

	resp, err := l.SubmitSatisfactionSurvey(req, 1, "participant")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "poster", resp.Feature)
	assert.Equal(suite.T(), &easeOfUse, resp.EaseOfUse)
	assert.Equal(suite.T(), "拖拽功能", resp.MostLiked)
}

func (suite *FeedbackLogicTestSuite) TestSubmitSatisfactionSurvey_Validation_EmptyFeature() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser")

	l := NewSubmitSatisfactionSurveyLogic(ctx, suite.svcCtx)
	req := &types.SubmitSatisfactionSurveyReq{}

	resp, err := l.SubmitSatisfactionSurvey(req, 1, "participant")

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "功能名称")
}

func (suite *FeedbackLogicTestSuite) TestSubmitSatisfactionSurvey_Validation_InvalidRating() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser")

	rating := 10
	l := NewSubmitSatisfactionSurveyLogic(ctx, suite.svcCtx)
	req := &types.SubmitSatisfactionSurveyReq{Feature: "poster", EaseOfUse: &rating}

	resp, err := l.SubmitSatisfactionSurvey(req, 1, "participant")

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "评分")
}

func (suite *FeedbackLogicTestSuite) TestListFAQ_Success() {
	ctx := context.Background()

	faq1 := &model.FAQItem{Category: "poster", Question: "如何创建海报？", Answer: "点击新建按钮即可创建", SortOrder: 1}
	faq2 := &model.FAQItem{Category: "payment", Question: "如何支付？", Answer: "支持微信和支付宝", SortOrder: 2}
	suite.db.Create(faq1)
	suite.db.Create(faq2)

	l := NewListFAQLogic(ctx, suite.svcCtx)
	req := &types.ListFAQReq{}

	resp, err := l.ListFAQ(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(2), resp.Total)
	assert.Len(suite.T(), resp.FAQs, 2)
}

func (suite *FeedbackLogicTestSuite) TestListFAQ_FilterByCategory() {
	ctx := context.Background()

	suite.db.Create(&model.FAQItem{Category: "poster", Question: "Q1", Answer: "A1", SortOrder: 1})
	suite.db.Create(&model.FAQItem{Category: "payment", Question: "Q2", Answer: "A2", SortOrder: 1})
	suite.db.Create(&model.FAQItem{Category: "poster", Question: "Q3", Answer: "A3", SortOrder: 2})

	l := NewListFAQLogic(ctx, suite.svcCtx)
	req := &types.ListFAQReq{Category: "poster"}

	resp, err := l.ListFAQ(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(2), resp.Total)
}

func (suite *FeedbackLogicTestSuite) TestListFAQ_SearchByKeyword() {
	ctx := context.Background()

	suite.db.Create(&model.FAQItem{Category: "poster", Question: "如何创建海报？", Answer: "点击新建", SortOrder: 1})
	suite.db.Create(&model.FAQItem{Category: "payment", Question: "支付问题", Answer: "联系客服", SortOrder: 1})
	suite.db.Create(&model.FAQItem{Category: "poster", Question: "删除海报", Answer: "点击删除按钮", SortOrder: 2})

	l := NewListFAQLogic(ctx, suite.svcCtx)
	req := &types.ListFAQReq{Keyword: "海报"}

	resp, err := l.ListFAQ(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(2), resp.Total)
}

func (suite *FeedbackLogicTestSuite) TestListFAQ_SortOrder() {
	ctx := context.Background()

	suite.db.Create(&model.FAQItem{Category: "poster", Question: "Q3", Answer: "A3", SortOrder: 3})
	suite.db.Create(&model.FAQItem{Category: "poster", Question: "Q1", Answer: "A1", SortOrder: 1})
	suite.db.Create(&model.FAQItem{Category: "poster", Question: "Q2", Answer: "A2", SortOrder: 2})

	l := NewListFAQLogic(ctx, suite.svcCtx)
	req := &types.ListFAQReq{}

	resp, err := l.ListFAQ(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(3), resp.Total)
	assert.Equal(suite.T(), 1, resp.FAQs[0].SortOrder)
	assert.Equal(suite.T(), 2, resp.FAQs[1].SortOrder)
	assert.Equal(suite.T(), 3, resp.FAQs[2].SortOrder)
}

func (suite *FeedbackLogicTestSuite) TestMarkFAQHelpful_Helpful() {
	ctx := context.Background()

	faq := &model.FAQItem{Category: "poster", Question: "Test Question", Answer: "Test Answer", SortOrder: 1, HelpfulCount: 5}
	suite.db.Create(faq)

	l := NewMarkFAQHelpfulLogic(ctx, suite.svcCtx)
	req := &types.MarkFAQHelpfulReq{Id: faq.ID, Type: "helpful"}

	resp, err := l.MarkFAQHelpful(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), faq.ID, resp.Id)
	assert.Equal(suite.T(), 6, resp.HelpfulCount)
}

func (suite *FeedbackLogicTestSuite) TestMarkFAQHelpful_NotHelpful() {
	ctx := context.Background()

	faq := &model.FAQItem{Category: "poster", Question: "Test Question", Answer: "Test Answer", SortOrder: 1, NotHelpfulCount: 2}
	suite.db.Create(faq)

	l := NewMarkFAQHelpfulLogic(ctx, suite.svcCtx)
	req := &types.MarkFAQHelpfulReq{Id: faq.ID, Type: "not_helpful"}

	resp, err := l.MarkFAQHelpful(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), faq.ID, resp.Id)
	assert.Equal(suite.T(), 3, resp.NotHelpfulCount)
}

func (suite *FeedbackLogicTestSuite) TestRecordFeatureUsage_Success() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser")
	campaignId := int64(100)
	durationMs := 1500

	l := NewRecordFeatureUsageLogic(ctx, suite.svcCtx)
	req := &types.RecordFeatureUsageReq{
		Feature:      "poster",
		Action:       "generate",
		CampaignId:   &campaignId,
		Success:      true,
		DurationMs:   &durationMs,
		ErrorMessage: "",
	}

	resp, err := l.RecordFeatureUsage(req, 1, "participant")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), 200, resp.Code)
	assert.Equal(suite.T(), "使用记录已保存", resp.Message)

	var stats []model.FeatureUsageStat
	suite.db.Where("feature = ? AND action = ?", req.Feature, req.Action).Find(&stats)
	assert.Len(suite.T(), stats, 1)
	assert.Equal(suite.T(), int64(1), stats[0].UserID)
	assert.Equal(suite.T(), "poster", stats[0].Feature)
	assert.Equal(suite.T(), "generate", stats[0].Action)
	assert.Equal(suite.T(), true, stats[0].Success)
	assert.Equal(suite.T(), &durationMs, stats[0].DurationMs)
}

func (suite *FeedbackLogicTestSuite) TestRecordFeatureUsage_Failure() {
	ctx := context.Background()
	suite.createTestUser(1, "testuser2")

	l := NewRecordFeatureUsageLogic(ctx, suite.svcCtx)
	req := &types.RecordFeatureUsageReq{
		Feature:      "poster2",
		Action:       "generate",
		Success:      false,
		DurationMs:   nil,
		ErrorMessage: "网络错误",
	}

	resp, err := l.RecordFeatureUsage(req, 1, "participant")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), 200, resp.Code)
	assert.Equal(suite.T(), "使用记录已保存", resp.Message)
}

func (suite *FeedbackLogicTestSuite) TestGetFeedbackStatistics_Success() {
	ctx := context.Background()
	suite.createTestUser(1, "admin")

	suite.db.Create(&model.UserFeedback{UserID: 1, Category: "poster", Title: "F1", Content: "C1", Status: "pending", Priority: "high", Rating: intPtr(5)})
	suite.db.Create(&model.UserFeedback{UserID: 1, Category: "payment", Title: "F2", Content: "C2", Status: "resolved", Priority: "medium", Rating: intPtr(4)})
	suite.db.Create(&model.UserFeedback{UserID: 1, Category: "poster", Title: "F3", Content: "C3", Status: "pending", Priority: "low", Rating: intPtr(3)})

	l := NewGetFeedbackStatisticsLogic(ctx, suite.svcCtx)
	req := &types.GetFeedbackStatisticsReq{}

	resp, err := l.GetFeedbackStatistics(req, "platform_admin")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(3), resp.TotalFeedbacks)
	assert.Equal(suite.T(), int64(2), resp.ByCategory["poster"])
	assert.Equal(suite.T(), int64(1), resp.ByCategory["payment"])
	assert.Equal(suite.T(), int64(2), resp.ByStatus["pending"])
	assert.Equal(suite.T(), int64(1), resp.ByStatus["resolved"])

	assert.Greater(suite.T(), resp.AverageRating, 0.0)
}

func (suite *FeedbackLogicTestSuite) TestGetFeedbackStatistics_ResolutionRate() {
	ctx := context.Background()
	suite.createTestUser(1, "admin")

	suite.db.Create(&model.UserFeedback{UserID: 1, Category: "poster", Title: "F1", Content: "C1", Status: "pending"})
	suite.db.Create(&model.UserFeedback{UserID: 1, Category: "poster", Title: "F2", Content: "C2", Status: "resolved"})
	suite.db.Create(&model.UserFeedback{UserID: 1, Category: "poster", Title: "F3", Content: "C3", Status: "resolved"})
	suite.db.Create(&model.UserFeedback{UserID: 1, Category: "poster", Title: "F4", Content: "C4", Status: "pending"})

	l := NewGetFeedbackStatisticsLogic(ctx, suite.svcCtx)
	req := &types.GetFeedbackStatisticsReq{}

	resp, err := l.GetFeedbackStatistics(req, "platform_admin")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(4), resp.TotalFeedbacks)
	assert.Equal(suite.T(), 0.5, resp.ResolutionRate)
}

func (suite *FeedbackLogicTestSuite) TestGetFeedbackStatistics_PermissionDenied() {
	ctx := context.Background()
	suite.createTestUser(1, "user")

	l := NewGetFeedbackStatisticsLogic(ctx, suite.svcCtx)
	req := &types.GetFeedbackStatisticsReq{}

	resp, err := l.GetFeedbackStatistics(req, "participant")

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "无权")
}

func (suite *FeedbackLogicTestSuite) TestGetFeedbackStatistics_DateRange() {
	ctx := context.Background()
	suite.createTestUser(1, "admin")

	oldTime := time.Now().AddDate(0, 0, -10)
	newTime := time.Now()

	suite.db.Create(&model.UserFeedback{UserID: 1, Category: "poster", Title: "Old Feedback", Content: "Old Content", Status: "pending", CreatedAt: oldTime})
	suite.db.Create(&model.UserFeedback{UserID: 1, Category: "poster", Title: "New Feedback", Content: "New Content", Status: "pending", CreatedAt: newTime})

	l := NewGetFeedbackStatisticsLogic(ctx, suite.svcCtx)
	req := &types.GetFeedbackStatisticsReq{StartDate: time.Now().AddDate(0, 0, -5).Format("2006-01-02")}

	resp, err := l.GetFeedbackStatistics(req, "platform_admin")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(1), resp.TotalFeedbacks)
}

func TestFeedbackLogicTestSuite(t *testing.T) {
	suite.Run(t, new(FeedbackLogicTestSuite))
}

func intPtr(i int) *int {
	return &i
}
