package order

import (
	"context"
	"fmt"
	"testing"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateOrderLogic_CampaignNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	logic := NewCreateOrderLogic(context.Background(), &svc.ServiceContext{DB: db})
	req := &types.CreateOrderReq{
		CampaignId: 9999,
		Phone:      "13800138000",
		FormData: map[string]string{
			"name": "张三",
		},
	}

	resp, err := logic.CreateOrder(req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "活动")
}

func TestCreateOrderLogic_CampaignEnded(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	campaign := &model.Campaign{
		Name:        "已结束活动",
		Description: "测试活动结束场景",
		FormFields:  `[{"type":"text","name":"name","label":"姓名","required":true}]`,
		RewardRule:  10,
		StartTime:   time.Now().Add(-48 * time.Hour),
		EndTime:     time.Now().Add(-24 * time.Hour),
		Status:      "active",
		BrandId:     1,
	}
	require.NoError(t, db.Create(campaign).Error)

	logic := NewCreateOrderLogic(context.Background(), &svc.ServiceContext{DB: db})
	req := &types.CreateOrderReq{
		CampaignId: campaign.Id,
		Phone:      "13800138000",
		FormData: map[string]string{
			"name": "张三",
		},
	}

	resp, err := logic.CreateOrder(req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "活动已结束")
}

func TestCreateOrderLogic_MissingRequiredField(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	campaign := createTestCampaign(t, db)
	logic := NewCreateOrderLogic(context.Background(), &svc.ServiceContext{DB: db})
	req := &types.CreateOrderReq{
		CampaignId: campaign.Id,
		Phone:      "13800138000",
		FormData:   map[string]string{},
	}

	resp, err := logic.CreateOrder(req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "必填字段")
}

func TestCreateOrderLogic_FormFieldValidationFailed(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	campaign := &model.Campaign{
		Name:        "邮箱验证活动",
		Description: "测试字段校验失败",
		FormFields:  `[{"type":"email","name":"email","label":"邮箱","required":true}]`,
		RewardRule:  10,
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     time.Now().Add(24 * time.Hour),
		Status:      "active",
		BrandId:     1,
	}
	require.NoError(t, db.Create(campaign).Error)

	logic := NewCreateOrderLogic(context.Background(), &svc.ServiceContext{DB: db})
	req := &types.CreateOrderReq{
		CampaignId: campaign.Id,
		Phone:      "13800138000",
		FormData: map[string]string{
			"email": "bad-email",
		},
	}

	resp, err := logic.CreateOrder(req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "字段 邮箱 验证失败")
}

func TestCreateOrderLogic_CampaignInactiveStatus(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	campaign := &model.Campaign{
		Name:        "暂停活动",
		Description: "测试活动状态非 active",
		FormFields:  `[{"type":"text","name":"name","label":"姓名","required":true}]`,
		RewardRule:  10,
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     time.Now().Add(24 * time.Hour),
		Status:      "paused",
		BrandId:     1,
	}
	require.NoError(t, db.Create(campaign).Error)

	logic := NewCreateOrderLogic(context.Background(), &svc.ServiceContext{DB: db})
	req := &types.CreateOrderReq{
		CampaignId: campaign.Id,
		Phone:      "13800138000",
		FormData: map[string]string{
			"name": "张三",
		},
	}

	resp, err := logic.CreateOrder(req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "活动未开始或已结束")
}

func TestValidateField_KeyTypes(t *testing.T) {
	tests := []struct {
		name    string
		field   model.FormField
		value   string
		hasErr  bool
		errText string
	}{
		{
			name:  "phone valid",
			field: model.FormField{Type: "phone"},
			value: "13800138000",
		},
		{
			name:    "phone invalid",
			field:   model.FormField{Type: "phone"},
			value:   "1380013",
			hasErr:  true,
			errText: "手机号格式不正确",
		},
		{
			name:  "email valid",
			field: model.FormField{Type: "email"},
			value: "user@example.com",
		},
		{
			name:    "email invalid",
			field:   model.FormField{Type: "email"},
			value:   "user@",
			hasErr:  true,
			errText: "邮箱格式不正确",
		},
		{
			name:  "number valid",
			field: model.FormField{Type: "number"},
			value: "100",
		},
		{
			name:    "number invalid",
			field:   model.FormField{Type: "number"},
			value:   "   ",
			hasErr:  true,
			errText: "数字不能为空",
		},
		{
			name:  "select valid",
			field: model.FormField{Type: "select", Options: []string{"A", "B"}},
			value: "A",
		},
		{
			name:    "select invalid",
			field:   model.FormField{Type: "select", Options: []string{"A", "B"}},
			value:   "C",
			hasErr:  true,
			errText: "请选择有效的选项",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := validateField(tc.value, tc.field)
			if tc.hasErr {
				require.Error(t, err)
				if tc.errText != "" {
					assert.Contains(t, err.Error(), tc.errText)
				}
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestValidateField_UnsupportedType(t *testing.T) {
	err := validateField("x", model.FormField{Type: "checkbox", Name: "agree", Label: "同意"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "不支持的字段类型")
}

func TestIsDuplicateOrderError_MySQL1062(t *testing.T) {
	err := &mysqlDriver.MySQLError{Number: 1062, Message: "Duplicate entry '1-13800138000' for key 'uk_orders_campaign_phone'"}
	assert.True(t, isDuplicateOrderError(err))

	wrapped := fmt.Errorf("db insert failed: %w", err)
	assert.True(t, isDuplicateOrderError(wrapped))
}

func TestIsDuplicateOrderError_SQLiteUniqueConstraint(t *testing.T) {
	err := fmt.Errorf("UNIQUE constraint failed: orders.campaign_id, orders.phone")
	assert.True(t, isDuplicateOrderError(err))
}

func TestIsDuplicateOrderError_NonDuplicate(t *testing.T) {
	err := fmt.Errorf("connection reset by peer")
	assert.False(t, isDuplicateOrderError(err))
}
