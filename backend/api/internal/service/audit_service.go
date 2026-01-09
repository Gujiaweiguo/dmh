package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"dmh/model"

	"gorm.io/gorm"
)

// AuditService 审计服务
type AuditService struct {
	db *gorm.DB
}

// NewAuditService 创建审计服务
func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{
		db: db,
	}
}

// AuditContext 审计上下文
type AuditContext struct {
	UserID    *int64
	Username  string
	ClientIP  string
	UserAgent string
}

// LogUserAction 记录用户操作日志
func (s *AuditService) LogUserAction(ctx *AuditContext, action, resource, resourceID string, details interface{}) error {
	detailsJSON := ""
	if details != nil {
		if detailsBytes, err := json.Marshal(details); err == nil {
			detailsJSON = string(detailsBytes)
		}
	}
	
	auditLog := &model.AuditLog{
		UserID:     ctx.UserID,
		Username:   ctx.Username,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Details:    detailsJSON,
		ClientIP:   ctx.ClientIP,
		UserAgent:  ctx.UserAgent,
		Status:     "success",
	}
	
	return s.db.Create(auditLog).Error
}

// LogFailedAction 记录失败的操作
func (s *AuditService) LogFailedAction(ctx *AuditContext, action, resource, resourceID string, errorMsg string) error {
	auditLog := &model.AuditLog{
		UserID:     ctx.UserID,
		Username:   ctx.Username,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		ClientIP:   ctx.ClientIP,
		UserAgent:  ctx.UserAgent,
		Status:     "failed",
		ErrorMsg:   errorMsg,
	}
	
	return s.db.Create(auditLog).Error
}

// LogLoginAttempt 记录登录尝试
func (s *AuditService) LogLoginAttempt(userID *int64, username, clientIP, userAgent string, success bool, failReason string) error {
	loginAttempt := &model.LoginAttempt{
		UserID:     userID,
		Username:   username,
		ClientIP:   clientIP,
		UserAgent:  userAgent,
		Success:    success,
		FailReason: failReason,
	}
	
	return s.db.Create(loginAttempt).Error
}

// LogSecurityEvent 记录安全事件
func (s *AuditService) LogSecurityEvent(eventType, severity string, userID *int64, username, clientIP, userAgent, description string, details interface{}) error {
	detailsJSON := ""
	if details != nil {
		if detailsBytes, err := json.Marshal(details); err == nil {
			detailsJSON = string(detailsBytes)
		}
	}
	
	securityEvent := &model.SecurityEvent{
		EventType:   eventType,
		Severity:    severity,
		UserID:      userID,
		Username:    username,
		ClientIP:    clientIP,
		UserAgent:   userAgent,
		Description: description,
		Details:     detailsJSON,
		Handled:     false,
	}
	
	return s.db.Create(securityEvent).Error
}

// GetAuditLogs 获取审计日志
func (s *AuditService) GetAuditLogs(page, pageSize int, filters map[string]interface{}) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64
	
	query := s.db.Model(&model.AuditLog{})
	
	// 应用过滤条件
	if userID, ok := filters["user_id"]; ok {
		query = query.Where("user_id = ?", userID)
	}
	if username, ok := filters["username"]; ok {
		query = query.Where("username LIKE ?", "%"+username.(string)+"%")
	}
	if action, ok := filters["action"]; ok {
		query = query.Where("action = ?", action)
	}
	if resource, ok := filters["resource"]; ok {
		query = query.Where("resource = ?", resource)
	}
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if startTime, ok := filters["start_time"]; ok {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime, ok := filters["end_time"]; ok {
		query = query.Where("created_at <= ?", endTime)
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&logs).Error
	
	return logs, total, err
}

// GetLoginAttempts 获取登录尝试记录
func (s *AuditService) GetLoginAttempts(page, pageSize int, filters map[string]interface{}) ([]model.LoginAttempt, int64, error) {
	var attempts []model.LoginAttempt
	var total int64
	
	query := s.db.Model(&model.LoginAttempt{})
	
	// 应用过滤条件
	if userID, ok := filters["user_id"]; ok {
		query = query.Where("user_id = ?", userID)
	}
	if username, ok := filters["username"]; ok {
		query = query.Where("username LIKE ?", "%"+username.(string)+"%")
	}
	if clientIP, ok := filters["client_ip"]; ok {
		query = query.Where("client_ip = ?", clientIP)
	}
	if success, ok := filters["success"]; ok {
		query = query.Where("success = ?", success)
	}
	if startTime, ok := filters["start_time"]; ok {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime, ok := filters["end_time"]; ok {
		query = query.Where("created_at <= ?", endTime)
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&attempts).Error
	
	return attempts, total, err
}

// GetSecurityEvents 获取安全事件
func (s *AuditService) GetSecurityEvents(page, pageSize int, filters map[string]interface{}) ([]model.SecurityEvent, int64, error) {
	var events []model.SecurityEvent
	var total int64
	
	query := s.db.Model(&model.SecurityEvent{})
	
	// 应用过滤条件
	if eventType, ok := filters["event_type"]; ok {
		query = query.Where("event_type = ?", eventType)
	}
	if severity, ok := filters["severity"]; ok {
		query = query.Where("severity = ?", severity)
	}
	if userID, ok := filters["user_id"]; ok {
		query = query.Where("user_id = ?", userID)
	}
	if handled, ok := filters["handled"]; ok {
		query = query.Where("handled = ?", handled)
	}
	if startTime, ok := filters["start_time"]; ok {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime, ok := filters["end_time"]; ok {
		query = query.Where("created_at <= ?", endTime)
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&events).Error
	
	return events, total, err
}

// HandleSecurityEvent 处理安全事件
func (s *AuditService) HandleSecurityEvent(eventID, handlerID int64, note string) error {
	updates := map[string]interface{}{
		"handled":    true,
		"handled_by": handlerID,
		"handled_at": time.Now(),
	}
	
	return s.db.Model(&model.SecurityEvent{}).
		Where("id = ?", eventID).
		Updates(updates).Error
}

// DetectSuspiciousActivity 检测可疑活动
func (s *AuditService) DetectSuspiciousActivity() error {
	// 检测频繁登录失败
	if err := s.detectFrequentLoginFailures(); err != nil {
		return err
	}
	
	// 检测异常IP登录
	if err := s.detectAbnormalIPLogin(); err != nil {
		return err
	}
	
	// 检测权限提升操作
	if err := s.detectPrivilegeEscalation(); err != nil {
		return err
	}
	
	return nil
}

// detectFrequentLoginFailures 检测频繁登录失败
func (s *AuditService) detectFrequentLoginFailures() error {
	// 查找最近1小时内登录失败超过10次的IP
	var results []struct {
		ClientIP string
		Count    int64
	}
	
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	err := s.db.Model(&model.LoginAttempt{}).
		Select("client_ip, COUNT(*) as count").
		Where("success = ? AND created_at >= ?", false, oneHourAgo).
		Group("client_ip").
		Having("COUNT(*) >= ?", 10).
		Find(&results).Error
	
	if err != nil {
		return err
	}
	
	// 记录安全事件
	for _, result := range results {
		s.LogSecurityEvent(
			"frequent_login_failures",
			"medium",
			nil,
			"",
			result.ClientIP,
			"",
			fmt.Sprintf("IP %s 在1小时内登录失败 %d 次", result.ClientIP, result.Count),
			map[string]interface{}{
				"client_ip":    result.ClientIP,
				"failure_count": result.Count,
				"time_window":  "1 hour",
			},
		)
	}
	
	return nil
}

// detectAbnormalIPLogin 检测异常IP登录
func (s *AuditService) detectAbnormalIPLogin() error {
	// 查找用户从新IP登录的情况
	var results []struct {
		UserID   int64
		Username string
		ClientIP string
	}
	
	// 查找最近24小时内的成功登录
	oneDayAgo := time.Now().Add(-24 * time.Hour)
	err := s.db.Raw(`
		SELECT DISTINCT la.user_id, la.username, la.client_ip
		FROM login_attempts la
		WHERE la.success = true 
		AND la.created_at >= ?
		AND la.user_id IS NOT NULL
		AND NOT EXISTS (
			SELECT 1 FROM login_attempts la2 
			WHERE la2.user_id = la.user_id 
			AND la2.client_ip = la.client_ip 
			AND la2.success = true 
			AND la2.created_at < ?
		)
	`, oneDayAgo, oneDayAgo).Scan(&results).Error
	
	if err != nil {
		return err
	}
	
	// 记录安全事件
	for _, result := range results {
		s.LogSecurityEvent(
			"abnormal_ip_login",
			"low",
			&result.UserID,
			result.Username,
			result.ClientIP,
			"",
			fmt.Sprintf("用户 %s 从新IP %s 登录", result.Username, result.ClientIP),
			map[string]interface{}{
				"user_id":   result.UserID,
				"username":  result.Username,
				"client_ip": result.ClientIP,
			},
		)
	}
	
	return nil
}

// detectPrivilegeEscalation 检测权限提升操作
func (s *AuditService) detectPrivilegeEscalation() error {
	// 查找最近1小时内的权限相关操作
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	var logs []model.AuditLog
	
	err := s.db.Where("action IN ? AND created_at >= ?", 
		[]string{"update_user_role", "assign_permission", "create_admin_user"}, 
		oneHourAgo).Find(&logs).Error
	
	if err != nil {
		return err
	}
	
	// 记录安全事件
	for _, log := range logs {
		s.LogSecurityEvent(
			"privilege_escalation",
			"high",
			log.UserID,
			log.Username,
			log.ClientIP,
			log.UserAgent,
			fmt.Sprintf("用户 %s 执行了权限提升操作: %s", log.Username, log.Action),
			map[string]interface{}{
				"action":      log.Action,
				"resource":    log.Resource,
				"resource_id": log.ResourceID,
				"details":     log.Details,
			},
		)
	}
	
	return nil
}

// GetAuditContextFromHTTP 从HTTP请求中获取审计上下文
func GetAuditContextFromHTTP(r *http.Request, userID *int64, username string) *AuditContext {
	return &AuditContext{
		UserID:    userID,
		Username:  username,
		ClientIP:  getClientIP(r),
		UserAgent: r.UserAgent(),
	}
}

// getClientIP 获取客户端真实IP
func getClientIP(r *http.Request) string {
	// 尝试从各种头部获取真实IP
	clientIP := r.Header.Get("X-Forwarded-For")
	if clientIP != "" {
		// X-Forwarded-For 可能包含多个IP，取第一个
		ips := strings.Split(clientIP, ",")
		return strings.TrimSpace(ips[0])
	}
	
	clientIP = r.Header.Get("X-Real-IP")
	if clientIP != "" {
		return clientIP
	}
	
	clientIP = r.Header.Get("X-Client-IP")
	if clientIP != "" {
		return clientIP
	}
	
	// 如果都没有，使用RemoteAddr
	return strings.Split(r.RemoteAddr, ":")[0]
}