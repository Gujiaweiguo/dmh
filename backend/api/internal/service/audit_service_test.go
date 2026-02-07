package service

import (
	"testing"

	"dmh/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AuditServiceTestSuite struct {
	suite.Suite
	db           *gorm.DB
	auditService *AuditService
}

func (suite *AuditServiceTestSuite) SetupSuite() {
	// 使用内存SQLite数据库进行测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	// 自动迁移表结构
	err = db.AutoMigrate(
		&model.AuditLog{},
		&model.LoginAttempt{},
		&model.SecurityEvent{},
		&model.User{},
	)
	suite.Require().NoError(err)

	suite.db = db
	suite.auditService = NewAuditService(db)
}

func (suite *AuditServiceTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func (suite *AuditServiceTestSuite) SetupTest() {
	// 清理测试数据
	suite.db.Exec("DELETE FROM audit_logs")
	suite.db.Exec("DELETE FROM login_attempts")
	suite.db.Exec("DELETE FROM security_events")
	suite.db.Exec("DELETE FROM users")
}

func (suite *AuditServiceTestSuite) TestLogUserAction() {
	userID := int64(1)
	ctx := &AuditContext{
		UserID:    &userID,
		Username:  "testuser",
		ClientIP:  "192.168.1.100",
		UserAgent: "Mozilla/5.0",
	}

	// 测试记录用户操作
	err := suite.auditService.LogUserAction(
		ctx,
		"create_user",
		"user",
		"123",
		map[string]interface{}{
			"username": "newuser",
			"role":     "participant",
		},
	)
	assert.NoError(suite.T(), err)

	// 验证日志已记录
	var log model.AuditLog
	err = suite.db.First(&log).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, *log.UserID)
	assert.Equal(suite.T(), "testuser", log.Username)
	assert.Equal(suite.T(), "create_user", log.Action)
	assert.Equal(suite.T(), "user", log.Resource)
	assert.Equal(suite.T(), "123", log.ResourceID)
	assert.Equal(suite.T(), "success", log.Status)
	assert.Equal(suite.T(), "192.168.1.100", log.ClientIP)
	assert.Contains(suite.T(), log.Details, "newuser")
}

func (suite *AuditServiceTestSuite) TestLogFailedAction() {
	userID := int64(1)
	ctx := &AuditContext{
		UserID:    &userID,
		Username:  "testuser",
		ClientIP:  "192.168.1.100",
		UserAgent: "Mozilla/5.0",
	}

	// 测试记录失败操作
	err := suite.auditService.LogFailedAction(
		ctx,
		"delete_user",
		"user",
		"456",
		"权限不足",
	)
	assert.NoError(suite.T(), err)

	// 验证日志已记录
	var log model.AuditLog
	err = suite.db.First(&log).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "failed", log.Status)
	assert.Equal(suite.T(), "权限不足", log.ErrorMsg)
}

func (suite *AuditServiceTestSuite) TestLogLoginAttempt() {
	userID := int64(1)

	// 测试记录成功登录
	err := suite.auditService.LogLoginAttempt(
		&userID,
		"testuser",
		"192.168.1.100",
		"Mozilla/5.0",
		true,
		"",
	)
	assert.NoError(suite.T(), err)

	// 测试记录失败登录
	err = suite.auditService.LogLoginAttempt(
		nil,
		"wronguser",
		"192.168.1.101",
		"Mozilla/5.0",
		false,
		"用户不存在",
	)
	assert.NoError(suite.T(), err)

	// 验证登录尝试记录
	var attempts []model.LoginAttempt
	err = suite.db.Find(&attempts).Error
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), attempts, 2)

	// 验证成功登录记录
	successAttempt := attempts[0]
	assert.Equal(suite.T(), userID, *successAttempt.UserID)
	assert.Equal(suite.T(), "testuser", successAttempt.Username)
	assert.True(suite.T(), successAttempt.Success)

	// 验证失败登录记录
	failAttempt := attempts[1]
	assert.Nil(suite.T(), failAttempt.UserID)
	assert.Equal(suite.T(), "wronguser", failAttempt.Username)
	assert.False(suite.T(), failAttempt.Success)
	assert.Equal(suite.T(), "用户不存在", failAttempt.FailReason)
}

func (suite *AuditServiceTestSuite) TestLogSecurityEvent() {
	userID := int64(1)

	// 测试记录安全事件
	err := suite.auditService.LogSecurityEvent(
		"frequent_login_failures",
		"medium",
		&userID,
		"testuser",
		"192.168.1.100",
		"Mozilla/5.0",
		"用户在短时间内多次登录失败",
		map[string]interface{}{
			"failure_count": 5,
			"time_window":   "1 hour",
		},
	)
	assert.NoError(suite.T(), err)

	// 验证安全事件已记录
	var event model.SecurityEvent
	err = suite.db.First(&event).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "frequent_login_failures", event.EventType)
	assert.Equal(suite.T(), "medium", event.Severity)
	assert.Equal(suite.T(), userID, *event.UserID)
	assert.Equal(suite.T(), "testuser", event.Username)
	assert.False(suite.T(), event.Handled)
	assert.Contains(suite.T(), event.Details, "failure_count")
}

func (suite *AuditServiceTestSuite) TestGetAuditLogs() {
	// 创建测试数据
	userID := int64(1)
	ctx := &AuditContext{
		UserID:    &userID,
		Username:  "testuser",
		ClientIP:  "192.168.1.100",
		UserAgent: "Mozilla/5.0",
	}

	// 记录多个操作
	for i := 0; i < 15; i++ {
		suite.auditService.LogUserAction(
			ctx,
			"test_action",
			"test_resource",
			string(rune(i)),
			nil,
		)
	}

	// 测试分页查询
	logs, total, err := suite.auditService.GetAuditLogs(1, 10, map[string]interface{}{})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(15), total)
	assert.Len(suite.T(), logs, 10)

	// 测试过滤查询
	logs, total, err = suite.auditService.GetAuditLogs(1, 10, map[string]interface{}{
		"username": "testuser",
		"action":   "test_action",
	})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(15), total)
	assert.Len(suite.T(), logs, 10)
}

func (suite *AuditServiceTestSuite) TestGetLoginAttempts() {
	// 创建测试数据
	userID := int64(1)
	for i := 0; i < 8; i++ {
		success := i%2 == 0
		suite.auditService.LogLoginAttempt(
			&userID,
			"testuser",
			"192.168.1.100",
			"Mozilla/5.0",
			success,
			"",
		)
	}

	// 测试分页查询
	attempts, total, err := suite.auditService.GetLoginAttempts(1, 5, map[string]interface{}{})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(8), total)
	assert.Len(suite.T(), attempts, 5)

	// 测试过滤查询 - 只查询成功的登录
	attempts, total, err = suite.auditService.GetLoginAttempts(1, 10, map[string]interface{}{
		"success": true,
	})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(4), total)
	assert.Len(suite.T(), attempts, 4)
}

func (suite *AuditServiceTestSuite) TestGetSecurityEvents() {
	// 创建测试数据
	userID := int64(1)
	eventTypes := []string{"frequent_login_failures", "abnormal_ip_login", "privilege_escalation"}
	severities := []string{"low", "medium", "high"}

	for i, eventType := range eventTypes {
		suite.auditService.LogSecurityEvent(
			eventType,
			severities[i],
			&userID,
			"testuser",
			"192.168.1.100",
			"Mozilla/5.0",
			"测试安全事件",
			nil,
		)
	}

	// 测试分页查询
	events, total, err := suite.auditService.GetSecurityEvents(1, 10, map[string]interface{}{})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), total)
	assert.Len(suite.T(), events, 3)

	// 测试过滤查询 - 只查询高严重程度事件
	events, total, err = suite.auditService.GetSecurityEvents(1, 10, map[string]interface{}{
		"severity": "high",
	})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), total)
	assert.Len(suite.T(), events, 1)
	assert.Equal(suite.T(), "privilege_escalation", events[0].EventType)
}

func (suite *AuditServiceTestSuite) TestHandleSecurityEvent() {
	// 创建测试安全事件
	userID := int64(1)
	handlerID := int64(2)

	err := suite.auditService.LogSecurityEvent(
		"test_event",
		"medium",
		&userID,
		"testuser",
		"192.168.1.100",
		"Mozilla/5.0",
		"测试安全事件",
		nil,
	)
	assert.NoError(suite.T(), err)

	// 获取事件ID
	var event model.SecurityEvent
	err = suite.db.First(&event).Error
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), event.Handled)

	// 处理安全事件
	err = suite.auditService.HandleSecurityEvent(event.ID, handlerID, "已处理")
	assert.NoError(suite.T(), err)

	// 验证事件已处理
	err = suite.db.First(&event, event.ID).Error
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), event.Handled)
	assert.Equal(suite.T(), handlerID, *event.HandledBy)
	assert.NotNil(suite.T(), event.HandledAt)
}

func (suite *AuditServiceTestSuite) TestDetectFrequentLoginFailures() {
	// 创建测试用户
	user := &model.User{
		Id:       1,
		Username: "testuser",
		Password: "hashedpassword",
		Phone:    "13800000001",
		Status:   "active",
	}
	err := suite.db.Create(user).Error
	assert.NoError(suite.T(), err)

	// 模拟频繁登录失败
	clientIP := "192.168.1.100"
	for i := 0; i < 12; i++ {
		suite.auditService.LogLoginAttempt(
			&user.Id,
			user.Username,
			clientIP,
			"Mozilla/5.0",
			false,
			"密码错误",
		)
	}

	// 执行可疑活动检测
	err = suite.auditService.DetectSuspiciousActivity()
	assert.NoError(suite.T(), err)

	// 验证安全事件已生成
	var events []model.SecurityEvent
	err = suite.db.Where("event_type = ?", "frequent_login_failures").Find(&events).Error
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), events, 1)
	assert.Equal(suite.T(), "medium", events[0].Severity)
	assert.Contains(suite.T(), events[0].Description, clientIP)
}

func TestAuditServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuditServiceTestSuite))
}
