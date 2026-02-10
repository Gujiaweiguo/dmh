package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"dmh/model"

	"gorm.io/gorm"
)

// SessionService 会话管理服务
type SessionService struct {
	db              *gorm.DB
	passwordService *PasswordService
}

// NewSessionService 创建会话管理服务
func NewSessionService(db *gorm.DB, passwordService *PasswordService) *SessionService {
	return &SessionService{
		db:              db,
		passwordService: passwordService,
	}
}

// CreateSession 创建用户会话
func (s *SessionService) CreateSession(userID int64, clientIP, userAgent string) (*model.UserSession, error) {
	policy, err := s.passwordService.GetPasswordPolicy()
	if err != nil {
		return nil, fmt.Errorf("获取会话策略失败: %v", err)
	}

	// 检查并发会话限制
	if err := s.checkConcurrentSessions(userID, policy.MaxConcurrentSessions); err != nil {
		return nil, err
	}

	// 生成会话ID
	sessionID, err := s.generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("生成会话ID失败: %v", err)
	}

	// 计算过期时间
	expiresAt := time.Now().Add(time.Duration(policy.SessionTimeout) * time.Minute)

	// 创建会话记录
	session := &model.UserSession{
		ID:           sessionID,
		UserID:       userID,
		ClientIP:     clientIP,
		UserAgent:    userAgent,
		LoginAt:      time.Now(),
		LastActiveAt: time.Now(),
		ExpiresAt:    expiresAt,
		Status:       "active",
	}

	if err := s.db.Create(session).Error; err != nil {
		return nil, fmt.Errorf("创建会话失败: %v", err)
	}

	return session, nil
}

// checkConcurrentSessions 检查并发会话限制
func (s *SessionService) checkConcurrentSessions(userID int64, maxSessions int) error {
	if maxSessions <= 0 {
		return nil // 无限制
	}

	// 清理过期会话
	s.cleanupExpiredSessions()

	// 统计当前活跃会话数
	var activeCount int64
	err := s.db.Model(&model.UserSession{}).
		Where("user_id = ? AND status = ? AND expires_at > ?", userID, "active", time.Now()).
		Count(&activeCount).Error

	if err != nil {
		return fmt.Errorf("检查会话数量失败: %v", err)
	}

	if activeCount >= int64(maxSessions) {
		// 强制下线最旧的会话
		if err := s.forceLogoutOldestSession(userID); err != nil {
			return fmt.Errorf("强制下线旧会话失败: %v", err)
		}
	}

	return nil
}

// forceLogoutOldestSession 强制下线最旧的会话
func (s *SessionService) forceLogoutOldestSession(userID int64) error {
	var oldestSession model.UserSession
	err := s.db.Where("user_id = ? AND status = ?", userID, "active").
		Order("login_at ASC").
		First(&oldestSession).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // 没有找到会话
		}
		return err
	}

	// 撤销最旧的会话
	return s.RevokeSession(oldestSession.ID, "concurrent_limit_exceeded")
}

// generateSessionID 生成会话ID
func (s *SessionService) generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GetSession 获取会话信息
func (s *SessionService) GetSession(sessionID string) (*model.UserSession, error) {
	var session model.UserSession
	err := s.db.Where("id = ?", sessionID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// ValidateSession 验证会话有效性
func (s *SessionService) ValidateSession(sessionID string) (*model.UserSession, error) {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("会话不存在")
	}

	// 检查会话状态
	if session.Status != "active" {
		return nil, fmt.Errorf("会话已失效")
	}

	// 检查会话是否过期
	if time.Now().After(session.ExpiresAt) {
		// 标记会话为过期
		s.db.Model(session).Update("status", "expired")
		return nil, fmt.Errorf("会话已过期")
	}

	return session, nil
}

// UpdateSessionActivity 更新会话活跃时间
func (s *SessionService) UpdateSessionActivity(sessionID string) error {
	policy, err := s.passwordService.GetPasswordPolicy()
	if err != nil {
		return err
	}

	now := time.Now()
	newExpiresAt := now.Add(time.Duration(policy.SessionTimeout) * time.Minute)

	updates := map[string]interface{}{
		"last_active_at": now,
		"expires_at":     newExpiresAt,
	}

	return s.db.Model(&model.UserSession{}).
		Where("id = ? AND status = ?", sessionID, "active").
		Updates(updates).Error
}

// RevokeSession 撤销会话
func (s *SessionService) RevokeSession(sessionID, reason string) error {
	updates := map[string]interface{}{
		"status": "revoked",
	}

	return s.db.Model(&model.UserSession{}).
		Where("id = ?", sessionID).
		Updates(updates).Error
}

// RevokeUserSessions 撤销用户的所有会话
func (s *SessionService) RevokeUserSessions(userID int64, excludeSessionID string) error {
	query := s.db.Model(&model.UserSession{}).Where("user_id = ? AND status = ?", userID, "active")

	if excludeSessionID != "" {
		query = query.Where("id != ?", excludeSessionID)
	}

	return query.Update("status", "revoked").Error
}

// GetUserSessions 获取用户的所有会话
func (s *SessionService) GetUserSessions(userID int64) ([]model.UserSession, error) {
	var sessions []model.UserSession
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&sessions).Error

	return sessions, err
}

// GetActiveSessions 获取活跃会话列表
func (s *SessionService) GetActiveSessions(page, pageSize int) ([]model.UserSession, int64, error) {
	var sessions []model.UserSession
	var total int64

	query := s.db.Model(&model.UserSession{}).Where("status = ? AND expires_at > ?", "active", time.Now())

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err := query.Order("last_active_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&sessions).Error

	return sessions, total, err
}

// cleanupExpiredSessions 清理过期会话
func (s *SessionService) cleanupExpiredSessions() error {
	return s.db.Model(&model.UserSession{}).
		Where("status = ? AND expires_at <= ?", "active", time.Now()).
		Update("status", "expired").Error
}

// CleanupOldSessions 清理旧会话记录
func (s *SessionService) CleanupOldSessions(retentionDays int) error {
	if retentionDays <= 0 {
		retentionDays = 30 // 默认保留30天
	}

	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	return s.db.Where("created_at < ?", cutoffTime).Delete(&model.UserSession{}).Error
}

// GetSessionStatistics 获取会话统计信息
func (s *SessionService) GetSessionStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 活跃会话数
	var activeCount int64
	s.db.Model(&model.UserSession{}).
		Where("status = ? AND expires_at > ?", "active", time.Now()).
		Count(&activeCount)
	stats["active_sessions"] = activeCount

	// 今日新建会话数
	today := time.Now().Truncate(24 * time.Hour)
	var todayCount int64
	s.db.Model(&model.UserSession{}).
		Where("created_at >= ?", today).
		Count(&todayCount)
	stats["today_sessions"] = todayCount

	// 过期会话数
	var expiredCount int64
	s.db.Model(&model.UserSession{}).
		Where("status = ? OR (status = ? AND expires_at <= ?)", "expired", "active", time.Now()).
		Count(&expiredCount)
	stats["expired_sessions"] = expiredCount

	// 撤销会话数
	var revokedCount int64
	s.db.Model(&model.UserSession{}).
		Where("status = ?", "revoked").
		Count(&revokedCount)
	stats["revoked_sessions"] = revokedCount

	// 平均会话时长（分钟）
	var avgDuration float64
	s.db.Model(&model.UserSession{}).
		Select("AVG(TIMESTAMPDIFF(MINUTE, login_at, COALESCE(updated_at, NOW())))").
		Where("status IN ?", []string{"expired", "revoked"}).
		Scan(&avgDuration)
	stats["avg_session_duration"] = avgDuration

	return stats, nil
}

// ForceLogoutUser 强制用户下线
func (s *SessionService) ForceLogoutUser(userID int64, reason string) error {
	// 撤销用户的所有活跃会话
	err := s.db.Model(&model.UserSession{}).
		Where("user_id = ? AND status = ?", userID, "active").
		Update("status", "revoked").Error

	if err != nil {
		return fmt.Errorf("强制下线用户失败: %v", err)
	}

	return nil
}

// IsUserOnline 检查用户是否在线
func (s *SessionService) IsUserOnline(userID int64) (bool, error) {
	var count int64
	err := s.db.Model(&model.UserSession{}).
		Where("user_id = ? AND status = ? AND expires_at > ?", userID, "active", time.Now()).
		Count(&count).Error

	return count > 0, err
}

// GetOnlineUsers 获取在线用户列表
func (s *SessionService) GetOnlineUsers() ([]int64, error) {
	var userIDs []int64
	err := s.db.Model(&model.UserSession{}).
		Select("DISTINCT user_id").
		Where("status = ? AND expires_at > ?", "active", time.Now()).
		Pluck("user_id", &userIDs).Error

	return userIDs, err
}
