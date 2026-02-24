package ai

import (
	"context"
	"testing"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AILogicTestSuite struct {
	suite.Suite
	svcCtx *svc.ServiceContext
}

func (suite *AILogicTestSuite) SetupSuite() {
	suite.svcCtx = &svc.ServiceContext{}
}

func (suite *AILogicTestSuite) TestGenerateCopywriting_DefaultStyle() {
	ctx := context.Background()
	l := NewGenerateCopywritingLogic(ctx, suite.svcCtx)
	req := &types.GenerateCopywritingReq{
		Topic: "新品发布",
	}

	resp, err := l.GenerateCopywriting(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Contains(suite.T(), resp.Content, "新品发布")
	assert.Contains(suite.T(), resp.Content, "📢")
}

func (suite *AILogicTestSuite) TestGenerateCopywriting_CasualStyle() {
	ctx := context.Background()
	l := NewGenerateCopywritingLogic(ctx, suite.svcCtx)
	req := &types.GenerateCopywritingReq{
		Topic: "夏季特惠",
		Style: "casual",
	}

	resp, err := l.GenerateCopywriting(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Contains(suite.T(), resp.Content, "夏季特惠")
	assert.Contains(suite.T(), resp.Content, "🎉")
}

func (suite *AILogicTestSuite) TestGenerateCopywriting_UrgentStyle() {
	ctx := context.Background()
	l := NewGenerateCopywritingLogic(ctx, suite.svcCtx)
	req := &types.GenerateCopywritingReq{
		Topic: "限时抢购",
		Style: "urgent",
	}

	resp, err := l.GenerateCopywriting(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Contains(suite.T(), resp.Content, "限时抢购")
	assert.Contains(suite.T(), resp.Content, "⚡")
}

func (suite *AILogicTestSuite) TestGenerateCopywriting_EmotionalStyle() {
	ctx := context.Background()
	l := NewGenerateCopywritingLogic(ctx, suite.svcCtx)
	req := &types.GenerateCopywritingReq{
		Topic: "家庭日",
		Style: "emotional",
	}

	resp, err := l.GenerateCopywriting(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Contains(suite.T(), resp.Content, "家庭日")
	assert.Contains(suite.T(), resp.Content, "❤️")
}

func (suite *AILogicTestSuite) TestGenerateCopywriting_ShortLength() {
	ctx := context.Background()
	l := NewGenerateCopywritingLogic(ctx, suite.svcCtx)
	req := &types.GenerateCopywritingReq{
		Topic:  "简短活动",
		Length: "short",
	}

	resp, err := l.GenerateCopywriting(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Contains(suite.T(), resp.Content, "简短活动")
	assert.Contains(suite.T(), resp.Content, "即刻参与")
}

func (suite *AILogicTestSuite) TestGenerateCopywriting_LongLength() {
	ctx := context.Background()
	l := NewGenerateCopywritingLogic(ctx, suite.svcCtx)
	req := &types.GenerateCopywritingReq{
		Topic:  "详细活动",
		Length: "long",
	}

	resp, err := l.GenerateCopywriting(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Contains(suite.T(), resp.Content, "详细活动")
	assert.Contains(suite.T(), resp.Content, "无论是品质还是服务")
}

func (suite *AILogicTestSuite) TestGenerateCopywriting_MediumLength() {
	ctx := context.Background()
	l := NewGenerateCopywritingLogic(ctx, suite.svcCtx)
	req := &types.GenerateCopywritingReq{
		Topic:  "中等活动",
		Length: "medium",
	}

	resp, err := l.GenerateCopywriting(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Contains(suite.T(), resp.Content, "中等活动")
	assert.Contains(suite.T(), resp.Content, "名额有限")
}

func (suite *AILogicTestSuite) TestGenerateCopywriting_CasualLong() {
	ctx := context.Background()
	l := NewGenerateCopywritingLogic(ctx, suite.svcCtx)
	req := &types.GenerateCopywritingReq{
		Topic:  "双11大促",
		Style:  "casual",
		Length: "long",
	}

	resp, err := l.GenerateCopywriting(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Contains(suite.T(), resp.Content, "双11大促")
	assert.Contains(suite.T(), resp.Content, "🎉")
	assert.Contains(suite.T(), resp.Content, "无论是品质还是服务")
}

func (suite *AILogicTestSuite) TestGenerateCopywriting_UrgentShort() {
	ctx := context.Background()
	l := NewGenerateCopywritingLogic(ctx, suite.svcCtx)
	req := &types.GenerateCopywritingReq{
		Topic:  "秒杀",
		Style:  "urgent",
		Length: "short",
	}

	resp, err := l.GenerateCopywriting(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Contains(suite.T(), resp.Content, "秒杀")
	assert.Contains(suite.T(), resp.Content, "⚡")
}

func TestAILogicTestSuite(t *testing.T) {
	suite.Run(t, new(AILogicTestSuite))
}
