package factory

import (
	"strings"
	"testing"
	"time"

	"dmh/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========== UserFactory Tests ==========

func TestUserFactory_Build(t *testing.T) {
	factory := NewUserFactory()
	user := factory.Build()

	assert.NotNil(t, user)
	assert.NotEmpty(t, user.Username)
	assert.NotEmpty(t, user.Phone)
	assert.NotEmpty(t, user.Email)
	assert.Equal(t, RoleParticipant, user.Role)
	assert.Equal(t, StatusActive, user.Status)
}

func TestUserFactory_BuildWith(t *testing.T) {
	factory := NewUserFactory()
	customUsername := "custom_user"
	customRole := RolePlatformAdmin

	user := factory.BuildWith(map[string]any{
		"username": customUsername,
		"role":     customRole,
	})

	assert.Equal(t, customUsername, user.Username)
	assert.Equal(t, customRole, user.Role)
}

func TestUserFactory_BuildWithSnakeCaseAndCamelCase(t *testing.T) {
	factory := NewUserFactory()

	// 测试 snake_case
	user1 := factory.BuildWith(map[string]any{
		"real_name": "snake_case_name",
	})
	assert.Equal(t, "snake_case_name", user1.RealName)

	// 测试 camelCase
	user2 := factory.BuildWith(map[string]any{
		"realName": "camelCaseName",
	})
	assert.Equal(t, "camelCaseName", user2.RealName)
}

func TestUserFactory_BuildList(t *testing.T) {
	factory := NewUserFactory()
	users := factory.BuildList(5)

	assert.Len(t, users, 5)
	for _, user := range users {
		assert.NotNil(t, user)
		assert.NotEmpty(t, user.Username)
	}
}

func TestUserFactory_BuildPlatformAdmin(t *testing.T) {
	factory := NewUserFactory()
	user := factory.BuildPlatformAdmin()

	assert.Equal(t, RolePlatformAdmin, user.Role)
	assert.Contains(t, user.Username, "admin")
}

func TestUserFactory_BuildBrandAdmin(t *testing.T) {
	factory := NewUserFactory()
	user := factory.BuildBrandAdmin()

	assert.Equal(t, RoleBrandAdmin, user.Role)
	assert.Contains(t, user.Username, "brand_admin")
}

func TestUserFactory_BuildLocked(t *testing.T) {
	factory := NewUserFactory()
	user := factory.BuildLocked()

	assert.Equal(t, StatusLocked, user.Status)
	assert.Equal(t, 5, user.LoginAttempts)
	assert.NotNil(t, user.LockedUntil)
}

// ========== BrandFactory Tests ==========

func TestBrandFactory_Build(t *testing.T) {
	factory := NewBrandFactory()
	brand := factory.Build()

	assert.NotNil(t, brand)
	assert.NotEmpty(t, brand.Name)
	assert.Equal(t, StatusActive, brand.Status)
}

func TestBrandFactory_BuildWith(t *testing.T) {
	factory := NewBrandFactory()
	customName := "Custom Brand"

	brand := factory.BuildWith(map[string]any{
		"name":   customName,
		"status": StatusDisabled,
	})

	assert.Equal(t, customName, brand.Name)
	assert.Equal(t, StatusDisabled, brand.Status)
}

// ========== CampaignFactory Tests ==========

func TestCampaignFactory_Build(t *testing.T) {
	factory := NewCampaignFactory()
	campaign := factory.Build()

	assert.NotNil(t, campaign)
	assert.NotEmpty(t, campaign.Name)
	assert.Equal(t, CampaignStatusActive, campaign.Status)
	assert.False(t, campaign.EnableDistribution)
}

func TestCampaignFactory_BuildWith(t *testing.T) {
	factory := NewCampaignFactory()
	brandId := int64(123)
	startTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	campaign := factory.BuildWith(map[string]any{
		"brand_id":   brandId,
		"name":       "Custom Campaign",
		"start_time": startTime,
	})

	assert.Equal(t, brandId, campaign.BrandId)
	assert.Equal(t, "Custom Campaign", campaign.Name)
	assert.Equal(t, startTime, campaign.StartTime)
}

func TestCampaignFactory_BuildWithDistribution(t *testing.T) {
	factory := NewCampaignFactory()
	brandId := int64(456)

	campaign := factory.BuildWithDistribution(brandId, 3)

	assert.Equal(t, brandId, campaign.BrandId)
	assert.True(t, campaign.EnableDistribution)
	assert.Equal(t, 3, campaign.DistributionLevel)
	assert.NotNil(t, campaign.DistributionRewards)
}

func TestCampaignFactory_BuildEnded(t *testing.T) {
	factory := NewCampaignFactory()
	brandId := int64(789)

	campaign := factory.BuildEnded(brandId)

	assert.Equal(t, brandId, campaign.BrandId)
	assert.Equal(t, CampaignStatusEnded, campaign.Status)
	assert.True(t, campaign.EndTime.Before(time.Now()))
}

// ========== OrderFactory Tests ==========

func TestOrderFactory_Build(t *testing.T) {
	factory := NewOrderFactory()
	order := factory.Build()

	assert.NotNil(t, order)
	assert.NotEmpty(t, order.Phone)
	assert.Equal(t, OrderStatusPending, order.Status)
	assert.Equal(t, PayStatusUnpaid, order.PayStatus)
}

func TestOrderFactory_BuildWith(t *testing.T) {
	factory := NewOrderFactory()
	campaignId := int64(999)
	customPhone := "13900139000"

	order := factory.BuildWith(map[string]any{
		"campaign_id": campaignId,
		"phone":       customPhone,
		"amount":      99.99,
	})

	assert.Equal(t, campaignId, order.CampaignId)
	assert.Equal(t, customPhone, order.Phone)
	assert.Equal(t, 99.99, order.Amount)
}

func TestOrderFactory_BuildPaid(t *testing.T) {
	factory := NewOrderFactory()
	campaignId := int64(111)

	order := factory.BuildPaid(campaignId)

	assert.Equal(t, campaignId, order.CampaignId)
	assert.Equal(t, OrderStatusPaid, order.Status)
	assert.Equal(t, PayStatusPaid, order.PayStatus)
	assert.Equal(t, 100.00, order.Amount)
	assert.NotNil(t, order.PaidAt)
}

func TestOrderFactory_BuildVerified(t *testing.T) {
	factory := NewOrderFactory()
	campaignId := int64(222)
	verifierId := int64(333)

	order := factory.BuildVerified(campaignId, verifierId)

	assert.Equal(t, campaignId, order.CampaignId)
	assert.Equal(t, "verified", order.VerificationStatus)
	assert.NotNil(t, order.VerifiedAt)
	assert.NotNil(t, order.VerifiedBy)
	assert.Equal(t, verifierId, *order.VerifiedBy)
	assert.NotEmpty(t, order.VerificationCode)
}

// ========== Error Tests ==========

func TestCampaignFactory_CreateWith_MissingBrandId(t *testing.T) {
	factory := NewCampaignFactory()
	// 传入 nil db，期望在验证 BrandId 时返回错误
	_, err := factory.CreateWith(nil, nil)
	assert.Equal(t, ErrBrandIdRequired, err)
}

func TestOrderFactory_CreateWith_MissingCampaignId(t *testing.T) {
	factory := NewOrderFactory()
	// 传入 nil db，期望在验证 CampaignId 时返回错误
	_, err := factory.CreateWith(nil, nil)
	assert.Equal(t, ErrCampaignIdRequired, err)
}

// ========== Integration Tests (需要数据库) ==========

// 这些测试需要真实的数据库连接，使用 build tag 或环境变量控制
// 运行方式: go test -tags=integration ./...

func TestUserFactory_Integration_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// 这里需要设置真实的数据库连接
	// 示例代码，实际测试时需要注入 db
	t.Skip("Requires database connection")
}

// ========== Benchmarks ==========

func BenchmarkUserFactory_Build(b *testing.B) {
	factory := NewUserFactory()
	for i := 0; i < b.N; i++ {
		factory.Build()
	}
}

func BenchmarkUserFactory_BuildWith(b *testing.B) {
	factory := NewUserFactory()
	overrides := map[string]any{
		"username": "bench_user",
		"role":     RoleParticipant,
	}
	for i := 0; i < b.N; i++ {
		factory.BuildWith(overrides)
	}
}

// ========== Helper function tests ==========

func TestBaseFactory_RandomSuffix(t *testing.T) {
	f := &BaseFactory{}
	suffixes := make(map[string]bool)

	// 验证随机性
	for i := 0; i < 100; i++ {
		s := f.RandomSuffix()
		assert.Len(t, s, 8) // 4 bytes = 8 hex chars
		suffixes[s] = true
	}

	// 验证足够随机（100次调用至少有90个不同值）
	assert.GreaterOrEqual(t, len(suffixes), 90)
}

func TestBaseFactory_RandomPhone(t *testing.T) {
	f := &BaseFactory{}
	phone := f.RandomPhone()

assert.Len(t, phone, 11)
	assert.True(t, strings.HasPrefix(phone, "138"))
}

func TestBaseFactory_RandomEmail(t *testing.T) {
	f := &BaseFactory{}
	email := f.RandomEmail()

	assert.Contains(t, email, "@example.com")
	assert.Contains(t, email, "test_")
}

func TestBaseFactory_TimeHelpers(t *testing.T) {
	f := &BaseFactory{}

	now := f.Now()
	future := f.FutureTime(time.Hour)
	past := f.PastTime(time.Hour)

	assert.True(t, future.After(now))
	assert.True(t, past.Before(now))
}

// ========== Model type assertion ==========

// 编译时类型检查，确保工厂返回正确的类型
func TestUserFactory_TypeAssertion(t *testing.T) {
	factory := NewUserFactory()
	var _ *model.User = factory.Build()
}

func TestBrandFactory_TypeAssertion(t *testing.T) {
	factory := NewBrandFactory()
	var _ *model.Brand = factory.Build()
}

func TestCampaignFactory_TypeAssertion(t *testing.T) {
	factory := NewCampaignFactory()
	var _ *model.Campaign = factory.Build()
}

func TestOrderFactory_TypeAssertion(t *testing.T) {
	factory := NewOrderFactory()
	var _ *model.Order = factory.Build()
}

// ========== Table-driven tests ==========

func TestUserFactory_BuildWith_TableDriven(t *testing.T) {
	factory := NewUserFactory()

	tests := []struct {
		name      string
		overrides map[string]any
		check     func(*testing.T, *model.User)
	}{
		{
			name:      "no overrides",
			overrides: nil,
			check: func(t *testing.T, u *model.User) {
				assert.Equal(t, RoleParticipant, u.Role)
				assert.Equal(t, StatusActive, u.Status)
			},
		},
		{
			name: "override role",
			overrides: map[string]any{
				"role": RolePlatformAdmin,
			},
			check: func(t *testing.T, u *model.User) {
				assert.Equal(t, RolePlatformAdmin, u.Role)
			},
		},
		{
			name: "override multiple fields",
			overrides: map[string]any{
				"username":       "multi_field_user",
				"role":           RoleBrandAdmin,
				"real_name":      "Multi Field User",
				"login_attempts": 3,
			},
			check: func(t *testing.T, u *model.User) {
				assert.Equal(t, "multi_field_user", u.Username)
				assert.Equal(t, RoleBrandAdmin, u.Role)
				assert.Equal(t, "Multi Field User", u.RealName)
				assert.Equal(t, 3, u.LoginAttempts)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := factory.BuildWith(tt.overrides)
			require.NotNil(t, user)
			tt.check(t, user)
		})
	}
}

func TestCampaignFactory_BuildWith_TableDriven(t *testing.T) {
	factory := NewCampaignFactory()

	tests := []struct {
		name      string
		overrides map[string]any
		check     func(*testing.T, *model.Campaign)
	}{
		{
			name:      "no overrides",
			overrides: nil,
			check: func(t *testing.T, c *model.Campaign) {
				assert.Equal(t, CampaignStatusActive, c.Status)
				assert.False(t, c.EnableDistribution)
			},
		},
		{
			name: "override brand_id and name",
			overrides: map[string]any{
				"brand_id": int64(123),
				"name":     "Custom Campaign",
			},
			check: func(t *testing.T, c *model.Campaign) {
				assert.Equal(t, int64(123), c.BrandId)
				assert.Equal(t, "Custom Campaign", c.Name)
			},
		},
		{
			name: "enable distribution",
			overrides: map[string]any{
				"brand_id":            int64(456),
				"enable_distribution": true,
				"distribution_level":  3,
			},
			check: func(t *testing.T, c *model.Campaign) {
				assert.True(t, c.EnableDistribution)
				assert.Equal(t, 3, c.DistributionLevel)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			campaign := factory.BuildWith(tt.overrides)
			require.NotNil(t, campaign)
			tt.check(t, campaign)
		})
	}
}
