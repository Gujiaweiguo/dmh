package service

import (
	"testing"
	"time"

	"dmh/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type PasswordServiceTestSuite struct {
	suite.Suite
	db              *gorm.DB
	passwordService *PasswordService
}

func (suite *PasswordServiceTestSuite) SetupSuite() {
	// 使用内存SQLite数据库进行测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	// 自动迁移表结构
	err = db.AutoMigrate(
		&model.PasswordPolicy{},
		&model.PasswordHistory{},
		&model.User{},
	)
	suite.Require().NoError(err)

	suite.db = db
	suite.passwordService = NewPasswordService(db)
}

func (suite *PasswordServiceTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func (suite *PasswordServiceTestSuite) SetupTest() {
	// 清理测试数据
	suite.db.Exec("DELETE FROM password_policies")
	suite.db.Exec("DELETE FROM password_histories")
	suite.db.Exec("DELETE FROM users")
}

func (suite *PasswordServiceTestSuite) TestGetPasswordPolicy() {
	// 测试获取默认密码策略
	policy, err := suite.passwordService.GetPasswordPolicy()
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), policy)
	assert.Equal(suite.T(), 8, policy.MinLength)
	assert.True(suite.T(), policy.RequireUppercase)
	assert.True(suite.T(), policy.RequireLowercase)
	assert.True(suite.T(), policy.RequireNumbers)
	assert.True(suite.T(), policy.RequireSpecialChars)

	// 测试获取自定义密码策略
	customPolicy := map[string]interface{}{
		"min_length":              10,
		"require_uppercase":       true,
		"require_lowercase":       true,
		"require_numbers":         true,
		"require_special_chars":   false,
		"max_age":                 60,
		"history_count":           3,
		"max_login_attempts":      3,
		"lockout_duration":        15,
		"session_timeout":         240,
		"max_concurrent_sessions": 2,
	}
	err = suite.db.Model(&model.PasswordPolicy{}).Create(customPolicy).Error
	assert.NoError(suite.T(), err)

	policy, err = suite.passwordService.GetPasswordPolicy()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 10, policy.MinLength)
	assert.False(suite.T(), policy.RequireSpecialChars)
}

func (suite *PasswordServiceTestSuite) TestValidatePassword() {
	// 测试有效密码
	err := suite.passwordService.ValidatePassword("Password123!", 0)
	assert.NoError(suite.T(), err)

	// 测试密码长度不足
	err = suite.passwordService.ValidatePassword("Pass1!", 0)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "密码长度不能少于")

	// 测试缺少大写字母
	err = suite.passwordService.ValidatePassword("password123!", 0)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "密码必须包含至少一个大写字母")

	// 测试缺少小写字母
	err = suite.passwordService.ValidatePassword("PASSWORD123!", 0)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "密码必须包含至少一个小写字母")

	// 测试缺少数字
	err = suite.passwordService.ValidatePassword("Password!", 0)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "密码必须包含至少一个数字")

	// 测试缺少特殊字符
	err = suite.passwordService.ValidatePassword("Password123", 0)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "密码必须包含至少一个特殊字符")
}

func (suite *PasswordServiceTestSuite) TestPasswordHistory() {
	userID := int64(1)

	// 创建测试用户
	user := &model.User{
		Id:       userID,
		Username: "testuser",
		Password: "hashedpassword",
		Phone:    "13800000001",
		Status:   "active",
	}
	err := suite.db.Create(user).Error
	assert.NoError(suite.T(), err)

	// 测试保存密码历史
	hashedPassword1, _ := suite.passwordService.HashPassword("Password123!")
	err = suite.passwordService.SavePasswordHistory(userID, hashedPassword1)
	assert.NoError(suite.T(), err)

	hashedPassword2, _ := suite.passwordService.HashPassword("NewPassword456!")
	err = suite.passwordService.SavePasswordHistory(userID, hashedPassword2)
	assert.NoError(suite.T(), err)

	// 测试密码历史检查 - 应该拒绝重复密码
	err = suite.passwordService.ValidatePassword("Password123!", userID)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "不能使用最近")

	// 测试新密码 - 应该通过
	err = suite.passwordService.ValidatePassword("AnotherPassword789!", userID)
	assert.NoError(suite.T(), err)
}

func (suite *PasswordServiceTestSuite) TestPasswordExpiration() {
	userID := int64(1)

	// 创建测试用户
	user := &model.User{
		Id:        userID,
		Username:  "testuser",
		Password:  "hashedpassword",
		Phone:     "13800000001",
		Status:    "active",
		CreatedAt: time.Now().AddDate(0, 0, -100), // 100天前创建
	}
	err := suite.db.Create(user).Error
	assert.NoError(suite.T(), err)

	// 测试密码过期检查
	expired, err := suite.passwordService.IsPasswordExpired(userID)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), expired) // 应该过期（默认90天）

	// 添加最近的密码历史
	hashedPassword, _ := suite.passwordService.HashPassword("RecentPassword123!")
	err = suite.passwordService.SavePasswordHistory(userID, hashedPassword)
	assert.NoError(suite.T(), err)

	// 再次检查 - 应该不过期
	expired, err = suite.passwordService.IsPasswordExpired(userID)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), expired)
}

func (suite *PasswordServiceTestSuite) TestPasswordStrength() {
	// 测试弱密码
	score := suite.passwordService.GeneratePasswordStrengthScore("123456")
	level := suite.passwordService.GetPasswordStrengthLevel(score)
	assert.True(suite.T(), score < 40)
	assert.Equal(suite.T(), "很弱", level)

	// 测试中等密码
	score = suite.passwordService.GeneratePasswordStrengthScore("Password123")
	level = suite.passwordService.GetPasswordStrengthLevel(score)
	assert.True(suite.T(), score >= 40 && score < 80)
	assert.Equal(suite.T(), "中等", level)

	// 测试强密码
	score = suite.passwordService.GeneratePasswordStrengthScore("MyVeryStr0ng!P@ssw0rd")
	level = suite.passwordService.GetPasswordStrengthLevel(score)
	assert.True(suite.T(), score >= 80)
	assert.Equal(suite.T(), "强", level)
}

func (suite *PasswordServiceTestSuite) TestHashAndVerifyPassword() {
	password := "TestPassword123!"

	// 测试密码加密
	hashedPassword, err := suite.passwordService.HashPassword(password)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), hashedPassword)
	assert.NotEqual(suite.T(), password, hashedPassword)

	// 测试密码验证
	err = suite.passwordService.VerifyPassword(password, hashedPassword)
	assert.NoError(suite.T(), err)

	// 测试错误密码验证
	err = suite.passwordService.VerifyPassword("WrongPassword", hashedPassword)
	assert.Error(suite.T(), err)
}

func (suite *PasswordServiceTestSuite) TestUpdatePasswordPolicy() {
	// 创建新的密码策略
	newPolicy := &model.PasswordPolicy{
		MinLength:             12,
		RequireUppercase:      true,
		RequireLowercase:      true,
		RequireNumbers:        true,
		RequireSpecialChars:   true,
		MaxAge:                30,
		HistoryCount:          10,
		MaxLoginAttempts:      3,
		LockoutDuration:       60,
		SessionTimeout:        120,
		MaxConcurrentSessions: 1,
	}

	err := suite.passwordService.UpdatePasswordPolicy(newPolicy)
	assert.NoError(suite.T(), err)

	// 验证策略已更新
	policy, err := suite.passwordService.GetPasswordPolicy()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 12, policy.MinLength)
	assert.Equal(suite.T(), 30, policy.MaxAge)
	assert.Equal(suite.T(), 10, policy.HistoryCount)
}

func TestPasswordServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PasswordServiceTestSuite))
}
