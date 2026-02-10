package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

// AuthMiddleware 权限验证中间件
type AuthMiddleware struct {
	logx.Logger
	jwtSecret string
}

// JWTClaims JWT声明结构
type JWTClaims struct {
	UserID   int64    `json:"userId"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	BrandIDs []int64  `json:"brandIds,omitempty"`
	jwt.RegisteredClaims
}

// AuthResponse 认证响应结构
type AuthResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		Logger:    logx.WithContext(context.Background()),
		jwtSecret: jwtSecret,
	}
}

// Handle 处理JWT认证和权限检查
func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 跳过不需要认证的路径
		if m.isPublicPath(r.URL.Path) {
			next(w, r)
			return
		}

		// 提取并验证JWT token
		token, err := m.extractToken(r)
		if err != nil {
			m.writeErrorResponse(w, http.StatusUnauthorized, "Token提取失败: "+err.Error())
			return
		}

		// 解析和验证token
		claims, err := m.validateToken(token)
		if err != nil {
			m.Logger.Errorf("Token验证失败: %v", err)
			m.writeErrorResponse(w, http.StatusUnauthorized, "Token验证失败: "+err.Error())
			return
		}

		// 检查token是否即将过期，如果是则返回刷新提示
		if m.shouldRefreshToken(claims) {
			w.Header().Set("X-Token-Refresh", "true")
		}

		// 将用户信息注入到context中
		ctx := m.injectUserContext(r.Context(), claims)
		r = r.WithContext(ctx)

		// 记录访问日志
		m.Logger.Infof("用户 %d (%s) 访问 %s %s, 角色: %v",
			claims.UserID, claims.Username, r.Method, r.URL.Path, claims.Roles)

		next(w, r)
	}
}

// extractToken 从请求中提取JWT token
func (m *AuthMiddleware) extractToken(r *http.Request) (string, error) {
	// 从Authorization header中提取
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1], nil
		}
		return "", fmt.Errorf("Authorization header格式错误")
	}

	// 从Cookie中提取（可选）
	cookie, err := r.Cookie("dmh_token")
	if err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}

	// 从查询参数中提取（可选，用于某些特殊场景）
	token := r.URL.Query().Get("token")
	if token != "" {
		return token, nil
	}

	return "", fmt.Errorf("未找到认证token")
}

// validateToken 验证JWT token并返回claims
func (m *AuthMiddleware) validateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return []byte(m.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token解析失败: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token无效")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("token claims类型错误")
	}

	// 验证token是否过期
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("token已过期")
	}

	// 验证必要字段
	if claims.UserID <= 0 {
		return nil, fmt.Errorf("token中用户ID无效")
	}

	if claims.Username == "" {
		return nil, fmt.Errorf("token中用户名为空")
	}

	if len(claims.Roles) == 0 {
		return nil, fmt.Errorf("token中角色信息为空")
	}

	return claims, nil
}

// shouldRefreshToken 检查是否应该刷新token
func (m *AuthMiddleware) shouldRefreshToken(claims *JWTClaims) bool {
	if claims.ExpiresAt == nil {
		return false
	}

	// 如果token在30分钟内过期，建议刷新
	refreshThreshold := time.Now().Add(30 * time.Minute)
	return claims.ExpiresAt.Time.Before(refreshThreshold)
}

// injectUserContext 将用户信息注入到context中
func (m *AuthMiddleware) injectUserContext(ctx context.Context, claims *JWTClaims) context.Context {
	ctx = context.WithValue(ctx, "userId", claims.UserID)
	ctx = context.WithValue(ctx, "username", claims.Username)
	ctx = context.WithValue(ctx, "roles", claims.Roles)
	ctx = context.WithValue(ctx, "brandIds", claims.BrandIDs)
	ctx = context.WithValue(ctx, "userClaims", claims)
	return ctx
}

// isPublicPath 检查是否为公开路径（不需要认证）
func (m *AuthMiddleware) isPublicPath(path string) bool {
	publicPaths := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/auth/refresh",
		"/api/v1/campaigns", // 匿名用户可以查看活动列表
		"/health",
		"/ping",
	}

	for _, publicPath := range publicPaths {
		if strings.HasPrefix(path, publicPath) {
			return true
		}
	}

	return false
}

// writeErrorResponse 写入错误响应
func (m *AuthMiddleware) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := AuthResponse{
		Code:    statusCode,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

// GenerateToken 生成JWT token（用于登录时调用）
func (m *AuthMiddleware) GenerateToken(userID int64, username string, roles []string, brandIDs []int64) (string, error) {
	now := time.Now()
	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		BrandIDs: brandIDs,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)), // 24小时过期
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "dmh-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.jwtSecret))
}

// RefreshToken 刷新JWT token
func (m *AuthMiddleware) RefreshToken(oldTokenString string) (string, error) {
	// 解析旧token（即使过期也要解析）
	token, err := jwt.ParseWithClaims(oldTokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.jwtSecret), nil
	})

	if err != nil {
		// 如果是过期错误，我们仍然可以刷新
		if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors&jwt.ValidationErrorExpired != 0 {
			// token过期，但我们可以从中提取claims来生成新token
		} else {
			return "", fmt.Errorf("token解析失败: %v", err)
		}
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return "", fmt.Errorf("token claims类型错误")
	}

	// 检查token是否在可刷新时间窗口内（例如过期后7天内可刷新）
	if claims.ExpiresAt != nil {
		maxRefreshTime := claims.ExpiresAt.Time.Add(7 * 24 * time.Hour)
		if time.Now().After(maxRefreshTime) {
			return "", fmt.Errorf("token过期时间过长，无法刷新")
		}
	}

	// 生成新token
	return m.GenerateToken(claims.UserID, claims.Username, claims.Roles, claims.BrandIDs)
}

// GetUserFromContext 从context中获取用户信息
func GetUserFromContext(ctx context.Context) (*JWTClaims, error) {
	claims, ok := ctx.Value("userClaims").(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("无法从context中获取用户信息")
	}
	return claims, nil
}

// GetUserIDFromContext 从context中获取用户ID
func GetUserIDFromContext(ctx context.Context) (int64, error) {
	switch value := ctx.Value("userId").(type) {
	case int64:
		return value, nil
	case int:
		return int64(value), nil
	case json.Number:
		parsed, err := value.Int64()
		if err != nil {
			return 0, fmt.Errorf("用户ID转换失败: %v", err)
		}
		return parsed, nil
	case float64:
		return int64(value), nil
	case string:
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("用户ID转换失败: %v", err)
		}
		return parsed, nil
	default:
		return 0, fmt.Errorf("无法从context中获取用户ID")
	}
}

// GetUserRolesFromContext 从context中获取用户角色
func GetUserRolesFromContext(ctx context.Context) ([]string, error) {
	switch roles := ctx.Value("roles").(type) {
	case []string:
		return roles, nil
	case []interface{}:
		roleStrings := make([]string, 0, len(roles))
		for _, role := range roles {
			if roleStr, ok := role.(string); ok {
				roleStrings = append(roleStrings, roleStr)
			}
		}
		if len(roleStrings) == 0 {
			return nil, fmt.Errorf("无法从context中获取用户角色")
		}
		return roleStrings, nil
	default:
		return nil, fmt.Errorf("无法从context中获取用户角色")
	}
}

// HasRole 检查用户是否有指定角色
func HasRole(ctx context.Context, role string) bool {
	roles, err := GetUserRolesFromContext(ctx)
	if err != nil {
		return false
	}

	for _, r := range roles {
		if r == role || r == "platform_admin" { // 平台管理员拥有所有角色权限
			return true
		}
	}
	return false
}

// IsPlatformAdmin 检查是否为平台管理员
func IsPlatformAdmin(ctx context.Context) bool {
	return HasRole(ctx, "platform_admin")
}

// GetUserBrandIDs 从context中获取用户管理的品牌ID列表（已废弃，保留兼容性）
func GetUserBrandIDs(ctx context.Context) ([]int64, error) {
	return []int64{}, nil // 返回空列表
}

// CanAccessBrand 检查用户是否可以访问指定品牌
func CanAccessBrand(ctx context.Context, brandID int64) bool {
	// 只有平台管理员可以访问所有品牌
	return IsPlatformAdmin(ctx)
}
