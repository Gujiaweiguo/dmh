package middleware

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// PermissionMiddleware 权限检查中间件
type PermissionMiddleware struct {
	logx.Logger
	db    *sql.DB
	cache *PermissionCache
}

// Permission 权限结构
type Permission struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

// UserPermissions 用户权限缓存结构
type UserPermissions struct {
	UserID      int64               `json:"userId"`
	Roles       []string            `json:"roles"`
	Permissions map[string]bool     `json:"permissions"` // permission_code -> bool
	BrandIDs    []int64             `json:"brandIds"`
	CachedAt    time.Time           `json:"cachedAt"`
}

// PermissionCache 权限缓存
type PermissionCache struct {
	mu    sync.RWMutex
	cache map[int64]*UserPermissions
	ttl   time.Duration
}

// PermissionResponse 权限检查响应
type PermissionResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewPermissionMiddleware(db *sql.DB) *PermissionMiddleware {
	return &PermissionMiddleware{
		Logger: logx.WithContext(context.Background()),
		db:     db,
		cache:  NewPermissionCache(30 * time.Minute), // 30分钟缓存
	}
}

func NewPermissionCache(ttl time.Duration) *PermissionCache {
	cache := &PermissionCache{
		cache: make(map[int64]*UserPermissions),
		ttl:   ttl,
	}
	
	// 启动清理goroutine
	go cache.cleanup()
	
	return cache
}

// Handle 处理权限检查
func (m *PermissionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 跳过不需要权限检查的路径
		if m.isPublicPath(r.URL.Path) {
			next(w, r)
			return
		}

		// 从context获取用户信息
		userID, err := GetUserIDFromContext(r.Context())
		if err != nil {
			m.writeErrorResponse(w, http.StatusUnauthorized, "用户信息获取失败")
			return
		}

		// 平台管理员拥有所有权限，直接放行
		if IsPlatformAdmin(r.Context()) {
			m.Logger.Infof("平台管理员 %d 访问 %s %s - 直接放行", userID, r.Method, r.URL.Path)
			next(w, r)
			return
		}

		// 提取资源和操作
		resource, action := m.extractResourceAction(r)
		permissionCode := fmt.Sprintf("%s:%s", resource, action)

		// 检查权限
		hasPermission, err := m.checkPermission(r.Context(), userID, permissionCode)
		if err != nil {
			m.Logger.Errorf("权限检查失败: %v", err)
			m.writeErrorResponse(w, http.StatusInternalServerError, "权限检查失败")
			return
		}

		if !hasPermission {
			m.Logger.Error(fmt.Sprintf("用户 %d 无权限访问 %s:%s", userID, resource, action))
			m.writeErrorResponse(w, http.StatusForbidden, fmt.Sprintf("无权限执行操作: %s", permissionCode))
			return
		}

		// 数据级权限检查
		if err := m.checkDataLevelPermission(r.Context(), r); err != nil {
			m.Logger.Error(fmt.Sprintf("用户 %d 数据级权限检查失败: %v", userID, err))
			m.writeErrorResponse(w, http.StatusForbidden, err.Error())
			return
		}

		m.Logger.Infof("用户 %d 权限检查通过: %s", userID, permissionCode)
		next(w, r)
	}
}

// checkPermission 检查用户是否有指定权限
func (m *PermissionMiddleware) checkPermission(ctx context.Context, userID int64, permissionCode string) (bool, error) {
	// 先从缓存获取
	userPerms, err := m.getUserPermissions(userID)
	if err != nil {
		return false, err
	}

	// 检查是否有该权限
	return userPerms.Permissions[permissionCode], nil
}

// getUserPermissions 获取用户权限（优先从缓存）
func (m *PermissionMiddleware) getUserPermissions(userID int64) (*UserPermissions, error) {
	// 先尝试从缓存获取
	if userPerms := m.cache.Get(userID); userPerms != nil {
		return userPerms, nil
	}

	// 缓存未命中，从数据库查询
	userPerms, err := m.loadUserPermissionsFromDB(userID)
	if err != nil {
		return nil, err
	}

	// 存入缓存
	m.cache.Set(userID, userPerms)
	
	return userPerms, nil
}

// loadUserPermissionsFromDB 从数据库加载用户权限
func (m *PermissionMiddleware) loadUserPermissionsFromDB(userID int64) (*UserPermissions, error) {
	userPerms := &UserPermissions{
		UserID:      userID,
		Permissions: make(map[string]bool),
		CachedAt:    time.Now(),
	}

	// 查询用户角色
	roleQuery := `
		SELECT r.code 
		FROM user_roles ur 
		JOIN roles r ON ur.role_id = r.id 
		WHERE ur.user_id = ?`
	
	rows, err := m.db.Query(roleQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户角色失败: %v", err)
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, fmt.Errorf("扫描角色失败: %v", err)
		}
		roles = append(roles, role)
	}
	userPerms.Roles = roles

	// 查询用户权限
	permQuery := `
		SELECT DISTINCT p.code 
		FROM user_roles ur 
		JOIN role_permissions rp ON ur.role_id = rp.role_id 
		JOIN permissions p ON rp.permission_id = p.id 
		WHERE ur.user_id = ?`
	
	rows, err = m.db.Query(permQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户权限失败: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var permCode string
		if err := rows.Scan(&permCode); err != nil {
			return nil, fmt.Errorf("扫描权限失败: %v", err)
		}
		userPerms.Permissions[permCode] = true
	}

	return userPerms, nil
}

// getUserBrandIDs 获取用户管理的品牌ID列表（已废弃，保留兼容性）
func (m *PermissionMiddleware) getUserBrandIDs(userID int64) ([]int64, error) {
	return []int64{}, nil
}

// checkDataLevelPermission 检查数据级权限（已简化，只有平台管理员可以访问所有数据）
func (m *PermissionMiddleware) checkDataLevelPermission(ctx context.Context, r *http.Request) error {
	// 只有平台管理员可以访问所有数据
	if !IsPlatformAdmin(ctx) {
		// 从URL中提取品牌ID
		brandID := m.extractBrandIDFromPath(r.URL.Path)
		if brandID != 0 {
			return fmt.Errorf("只有平台管理员可以访问品牌数据")
		}
	}
	return nil
}

// extractResourceAction 从请求中提取资源和操作
func (m *PermissionMiddleware) extractResourceAction(r *http.Request) (string, string) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/")
	parts := strings.Split(path, "/")
	
	var resource, action string
	
	if len(parts) > 0 {
		resource = parts[0]
	}
	
	// 根据HTTP方法和路径模式确定操作
	switch r.Method {
	case http.MethodGet:
		action = "read"
	case http.MethodPost:
		// 特殊处理一些POST操作
		if strings.Contains(path, "/approve") {
			action = "approve"
		} else if strings.Contains(path, "/reset-password") {
			action = "reset-password"
		} else if strings.Contains(path, "/status") {
			action = "status"
		} else {
			action = "create"
		}
	case http.MethodPut, http.MethodPatch:
		action = "update"
	case http.MethodDelete:
		action = "delete"
	default:
		action = "read"
	}
	
	return resource, action
}

// extractBrandIDFromPath 从URL路径中提取品牌ID
func (m *PermissionMiddleware) extractBrandIDFromPath(path string) int64 {
	// 匹配 /api/v1/brands/{id} 或 /api/v1/brands/{id}/xxx 模式
	if strings.Contains(path, "/brands/") {
		parts := strings.Split(path, "/brands/")
		if len(parts) > 1 {
			idPart := strings.Split(parts[1], "/")[0]
			if id, err := strconv.ParseInt(idPart, 10, 64); err == nil {
				return id
			}
		}
	}
	
	// 匹配其他可能包含品牌ID的路径
	// 例如 /api/v1/campaigns?brandId=1
	// 这种情况在业务逻辑中处理
	
	return 0
}

// isPublicPath 检查是否为公开路径
func (m *PermissionMiddleware) isPublicPath(path string) bool {
	publicPaths := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/register", 
		"/api/v1/auth/refresh",
		"/health",
		"/ping",
	}

	for _, publicPath := range publicPaths {
		if strings.HasPrefix(path, publicPath) {
			return true
		}
	}

	// 匿名用户可以查看活动列表（GET请求）
	if strings.HasPrefix(path, "/api/v1/campaigns") && strings.Contains(path, "GET") {
		return true
	}

	return false
}

// writeErrorResponse 写入错误响应
func (m *PermissionMiddleware) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := PermissionResponse{
		Code:    statusCode,
		Message: message,
	}
	
	json.NewEncoder(w).Encode(response)
}

// ClearUserCache 清除用户权限缓存（用户权限变更时调用）
func (m *PermissionMiddleware) ClearUserCache(userID int64) {
	m.cache.Delete(userID)
	m.Logger.Infof("已清除用户 %d 的权限缓存", userID)
}

// ClearAllCache 清除所有权限缓存
func (m *PermissionMiddleware) ClearAllCache() {
	m.cache.Clear()
	m.Logger.Info("已清除所有权限缓存")
}

// PermissionCache 方法实现

// Get 从缓存获取用户权限
func (c *PermissionCache) Get(userID int64) *UserPermissions {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	userPerms, exists := c.cache[userID]
	if !exists {
		return nil
	}
	
	// 检查是否过期
	if time.Since(userPerms.CachedAt) > c.ttl {
		delete(c.cache, userID)
		return nil
	}
	
	return userPerms
}

// Set 设置用户权限到缓存
func (c *PermissionCache) Set(userID int64, userPerms *UserPermissions) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.cache[userID] = userPerms
}

// Delete 删除用户权限缓存
func (c *PermissionCache) Delete(userID int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	delete(c.cache, userID)
}

// Clear 清除所有缓存
func (c *PermissionCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.cache = make(map[int64]*UserPermissions)
}

// cleanup 定期清理过期缓存
func (c *PermissionCache) cleanup() {
	ticker := time.NewTicker(10 * time.Minute) // 每10分钟清理一次
	defer ticker.Stop()
	
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for userID, userPerms := range c.cache {
			if now.Sub(userPerms.CachedAt) > c.ttl {
				delete(c.cache, userID)
			}
		}
		c.mu.Unlock()
	}
}

// ValidatePermission 验证权限格式
func ValidatePermission(permissionCode string) error {
	parts := strings.Split(permissionCode, ":")
	if len(parts) != 2 {
		return fmt.Errorf("权限格式错误，应为 resource:action 格式")
	}
	
	resource, action := parts[0], parts[1]
	if resource == "" || action == "" {
		return fmt.Errorf("资源和操作不能为空")
	}
	
	return nil
}