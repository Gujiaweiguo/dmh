package service

import (
	"testing"
	"time"

	"dmh/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SessionServiceTestSuite struct {
	suite.Suite
	db              *gorm.DB
	sessionService  *SessionService
	passwordService *PasswordService
}

func (suite *SessionServiceTestSuite) SetupSuite() {
	// 使用内存SQLite数据库进行测试
	db, err := gorm.Open(mysql.Open("root:Admin168@tcp(127.0.0.1:3306)/dmh_test?charset=utf8mb4&parseTime=true&loc=Local"), &gorm.Config{})
	suite.Require().NoError(err)

	// 自动迁移表结构
	err = db.AutoMigrate(
		&model.UserSession{},
		&model.PasswordPolicy{},
		&model.User{},
	)
	suite.Require().NoError(err)

	suite.db = db
	suite.passwordService = NewPasswordService(db)
	suite.sessionService = NewSessionService(db, suite.passwordService)
}

func (suite *SessionServiceTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func (suite *SessionServiceTestSuite) SetupTest() {
	suite.Require().NoError(suite.db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error)
	suite.Require().NoError(suite.db.Exec("TRUNCATE TABLE user_sessions").Error)
	suite.Require().NoError(suite.db.Exec("TRUNCATE TABLE password_histories").Error)
	suite.Require().NoError(suite.db.Exec("TRUNCATE TABLE password_policies").Error)
	suite.Require().NoError(suite.db.Exec("TRUNCATE TABLE user_feedback").Error)
	suite.Require().NoError(suite.db.Exec("TRUNCATE TABLE users").Error)
	suite.Require().NoError(suite.db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error)

	// 创建默认密码策略
	policy := &model.PasswordPolicy{
		MinLength:             8,
		RequireUppercase:      true,
		RequireLowercase:      true,
		RequireNumbers:        true,
		RequireSpecialChars:   true,
		MaxAge:                90,
		HistoryCount:          5,
		MaxLoginAttempts:      5,
		LockoutDuration:       30,
		SessionTimeout:        480,
		MaxConcurrentSessions: 3,
	}
	suite.Require().NoError(suite.db.Create(policy).Error)
}

func (suite *SessionServiceTestSuite) TestCreateSession() {
	userID := int64(1)
	clientIP := "192.168.1.100"
	userAgent := "Mozilla/5.0"

	// 测试创建会话
	session, err := suite.sessionService.CreateSession(userID, clientIP, userAgent)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), session)
	assert.Equal(suite.T(), userID, session.UserID)
	assert.Equal(suite.T(), clientIP, session.ClientIP)
	assert.Equal(suite.T(), userAgent, session.UserAgent)
	assert.Equal(suite.T(), "active", session.Status)
	assert.NotEmpty(suite.T(), session.ID)
	assert.True(suite.T(), session.ExpiresAt.After(time.Now()))
}

func (suite *SessionServiceTestSuite) TestConcurrentSessionLimit() {
	userID := int64(1)
	clientIP := "192.168.1.100"
	userAgent := "Mozilla/5.0"

	// 创建3个会话（达到限制）
	for i := 0; i < 3; i++ {
		session, err := suite.sessionService.CreateSession(userID, clientIP, userAgent)
		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), session)
	}

	// 验证有3个活跃会话
	sessions, err := suite.sessionService.GetUserSessions(userID)
	assert.NoError(suite.T(), err)
	activeCount := 0
	for _, s := range sessions {
		if s.Status == "active" {
			activeCount++
		}
	}
	assert.Equal(suite.T(), 3, activeCount)

	// 创建第4个会话，应该强制下线最旧的会话
	session, err := suite.sessionService.CreateSession(userID, clientIP, userAgent)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), session)

	// 验证仍然只有3个活跃会话
	sessions, err = suite.sessionService.GetUserSessions(userID)
	assert.NoError(suite.T(), err)
	activeCount = 0
	for _, s := range sessions {
		if s.Status == "active" {
			activeCount++
		}
	}
	assert.Equal(suite.T(), 3, activeCount)
}

func (suite *SessionServiceTestSuite) TestValidateSession() {
	userID := int64(1)
	clientIP := "192.168.1.100"
	userAgent := "Mozilla/5.0"

	// 创建会话
	session, err := suite.sessionService.CreateSession(userID, clientIP, userAgent)
	assert.NoError(suite.T(), err)

	// 测试验证有效会话
	validSession, err := suite.sessionService.ValidateSession(session.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), session.ID, validSession.ID)

	// 测试验证不存在的会话
	_, err = suite.sessionService.ValidateSession("nonexistent")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "会话不存在")

	// 撤销会话
	err = suite.sessionService.RevokeSession(session.ID, "test")
	assert.NoError(suite.T(), err)

	// 测试验证已撤销的会话
	_, err = suite.sessionService.ValidateSession(session.ID)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "会话已失效")
}

func (suite *SessionServiceTestSuite) TestSessionExpiration() {
	userID := int64(1)
	clientIP := "192.168.1.100"
	userAgent := "Mozilla/5.0"

	// 创建会话
	session, err := suite.sessionService.CreateSession(userID, clientIP, userAgent)
	assert.NoError(suite.T(), err)

	// 手动设置会话过期
	expiredTime := time.Now().Add(-1 * time.Hour)
	err = suite.db.Model(&model.UserSession{}).
		Where("id = ?", session.ID).
		Update("expires_at", expiredTime).Error
	assert.NoError(suite.T(), err)

	// 测试验证过期会话
	_, err = suite.sessionService.ValidateSession(session.ID)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "会话已过期")

	// 验证会话状态已更新为过期
	var updatedSession model.UserSession
	err = suite.db.Where("id = ?", session.ID).First(&updatedSession).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "expired", updatedSession.Status)
}

func (suite *SessionServiceTestSuite) TestUpdateSessionActivity() {
	userID := int64(1)
	clientIP := "192.168.1.100"
	userAgent := "Mozilla/5.0"

	// 创建会话
	session, err := suite.sessionService.CreateSession(userID, clientIP, userAgent)
	assert.NoError(suite.T(), err)

	originalExpiresAt := session.ExpiresAt
	originalLastActiveAt := session.LastActiveAt

	// 等待一秒确保时间差异
	time.Sleep(1 * time.Second)

	// 更新会话活跃时间
	err = suite.sessionService.UpdateSessionActivity(session.ID)
	assert.NoError(suite.T(), err)

	// 验证会话时间已更新
	var updatedSession model.UserSession
	err = suite.db.Where("id = ?", session.ID).First(&updatedSession).Error
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), updatedSession.LastActiveAt.After(originalLastActiveAt))
	assert.True(suite.T(), updatedSession.ExpiresAt.After(originalExpiresAt))
}

func (suite *SessionServiceTestSuite) TestRevokeSession() {
	userID := int64(1)
	clientIP := "192.168.1.100"
	userAgent := "Mozilla/5.0"

	// 创建会话
	session, err := suite.sessionService.CreateSession(userID, clientIP, userAgent)
	assert.NoError(suite.T(), err)

	// 撤销会话
	err = suite.sessionService.RevokeSession(session.ID, "admin_revoked")
	assert.NoError(suite.T(), err)

	// 验证会话状态已更新
	var updatedSession model.UserSession
	err = suite.db.Where("id = ?", session.ID).First(&updatedSession).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "revoked", updatedSession.Status)
}

func (suite *SessionServiceTestSuite) TestRevokeUserSessions() {
	userID := int64(1)
	clientIP := "192.168.1.100"
	userAgent := "Mozilla/5.0"

	// 创建多个会话
	var sessions []*model.UserSession
	for i := 0; i < 3; i++ {
		session, err := suite.sessionService.CreateSession(userID, clientIP, userAgent)
		assert.NoError(suite.T(), err)
		sessions = append(sessions, session)
	}

	// 撤销用户的所有会话，除了第一个
	err := suite.sessionService.RevokeUserSessions(userID, sessions[0].ID)
	assert.NoError(suite.T(), err)

	// 验证只有第一个会话仍然活跃
	userSessions, err := suite.sessionService.GetUserSessions(userID)
	assert.NoError(suite.T(), err)
	activeCount := 0
	for _, s := range userSessions {
		if s.Status == "active" {
			activeCount++
			assert.Equal(suite.T(), sessions[0].ID, s.ID)
		}
	}
	assert.Equal(suite.T(), 1, activeCount)
}

func (suite *SessionServiceTestSuite) TestGetActiveSessions() {
	// 创建多个用户的会话
	for userID := int64(1); userID <= 5; userID++ {
		for i := 0; i < 2; i++ {
			_, err := suite.sessionService.CreateSession(
				userID,
				"192.168.1.100",
				"Mozilla/5.0",
			)
			assert.NoError(suite.T(), err)
		}
	}

	// 测试分页查询活跃会话
	sessions, total, err := suite.sessionService.GetActiveSessions(1, 5)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(10), total)
	assert.Len(suite.T(), sessions, 5)

	// 验证所有返回的会话都是活跃的
	for _, session := range sessions {
		assert.Equal(suite.T(), "active", session.Status)
		assert.True(suite.T(), session.ExpiresAt.After(time.Now()))
	}
}

func (suite *SessionServiceTestSuite) TestForceLogoutUser() {
	userID := int64(1)
	clientIP := "192.168.1.100"
	userAgent := "Mozilla/5.0"

	// 创建多个会话
	for i := 0; i < 3; i++ {
		_, err := suite.sessionService.CreateSession(userID, clientIP, userAgent)
		assert.NoError(suite.T(), err)
	}

	// 强制用户下线
	err := suite.sessionService.ForceLogoutUser(userID, "admin_action")
	assert.NoError(suite.T(), err)

	// 验证用户的所有会话都被撤销
	sessions, err := suite.sessionService.GetUserSessions(userID)
	assert.NoError(suite.T(), err)
	for _, session := range sessions {
		assert.Equal(suite.T(), "revoked", session.Status)
	}
}

func (suite *SessionServiceTestSuite) TestIsUserOnline() {
	userID := int64(1)
	clientIP := "192.168.1.100"
	userAgent := "Mozilla/5.0"

	// 用户未登录时应该返回false
	online, err := suite.sessionService.IsUserOnline(userID)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), online)

	// 创建会话后应该返回true
	_, err = suite.sessionService.CreateSession(userID, clientIP, userAgent)
	assert.NoError(suite.T(), err)

	online, err = suite.sessionService.IsUserOnline(userID)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), online)

	// 强制下线后应该返回false
	err = suite.sessionService.ForceLogoutUser(userID, "test")
	assert.NoError(suite.T(), err)

	online, err = suite.sessionService.IsUserOnline(userID)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), online)
}

func (suite *SessionServiceTestSuite) TestGetOnlineUsers() {
	// 创建多个用户的会话
	userIDs := []int64{1, 2, 3, 4, 5}
	for _, userID := range userIDs {
		_, err := suite.sessionService.CreateSession(
			userID,
			"192.168.1.100",
			"Mozilla/5.0",
		)
		assert.NoError(suite.T(), err)
	}

	// 强制部分用户下线
	err := suite.sessionService.ForceLogoutUser(userIDs[0], "test")
	assert.NoError(suite.T(), err)
	err = suite.sessionService.ForceLogoutUser(userIDs[1], "test")
	assert.NoError(suite.T(), err)

	// 获取在线用户列表
	onlineUsers, err := suite.sessionService.GetOnlineUsers()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), onlineUsers, 3)

	// 验证在线用户ID
	expectedOnlineUsers := []int64{3, 4, 5}
	for _, expectedUserID := range expectedOnlineUsers {
		assert.Contains(suite.T(), onlineUsers, expectedUserID)
	}
}

func (suite *SessionServiceTestSuite) TestGetSessionStatistics() {
	userID := int64(1)
	clientIP := "192.168.1.100"
	userAgent := "Mozilla/5.0"

	s1, err := suite.sessionService.CreateSession(userID, clientIP, userAgent)
	assert.NoError(suite.T(), err)

	s2, err := suite.sessionService.CreateSession(userID, clientIP, userAgent)
	assert.NoError(suite.T(), err)

	err = suite.sessionService.RevokeSession(s2.ID, "unit-test")
	assert.NoError(suite.T(), err)

	err = suite.db.Model(&model.UserSession{}).
		Where("id = ?", s1.ID).
		Update("expires_at", time.Now().Add(-1*time.Hour)).Error
	assert.NoError(suite.T(), err)

	stats, err := suite.sessionService.GetSessionStatistics()
	assert.NoError(suite.T(), err)

	activeSessions, ok := stats["active_sessions"].(int64)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), int64(0), activeSessions)

	todaySessions, ok := stats["today_sessions"].(int64)
	assert.True(suite.T(), ok)
	assert.GreaterOrEqual(suite.T(), todaySessions, int64(2))

	expiredSessions, ok := stats["expired_sessions"].(int64)
	assert.True(suite.T(), ok)
	assert.GreaterOrEqual(suite.T(), expiredSessions, int64(1))

	revokedSessions, ok := stats["revoked_sessions"].(int64)
	assert.True(suite.T(), ok)
	assert.GreaterOrEqual(suite.T(), revokedSessions, int64(1))

	_, ok = stats["avg_session_duration"].(float64)
	assert.True(suite.T(), ok)
}

func (suite *SessionServiceTestSuite) TestCleanupOldSessions() {
	userID := int64(1)
	clientIP := "192.168.1.100"
	userAgent := "Mozilla/5.0"

	// 创建会话
	session, err := suite.sessionService.CreateSession(userID, clientIP, userAgent)
	assert.NoError(suite.T(), err)

	// 手动设置会话创建时间为31天前
	oldTime := time.Now().AddDate(0, 0, -31)
	err = suite.db.Model(&model.UserSession{}).
		Where("id = ?", session.ID).
		Update("created_at", oldTime).Error
	assert.NoError(suite.T(), err)

	// 清理30天前的会话
	err = suite.sessionService.CleanupOldSessions(30)
	assert.NoError(suite.T(), err)

	// 验证旧会话已被删除
	var count int64
	err = suite.db.Model(&model.UserSession{}).Where("id = ?", session.ID).Count(&count).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(0), count)
}

func TestSessionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SessionServiceTestSuite))
}
