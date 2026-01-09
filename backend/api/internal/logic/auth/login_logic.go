// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dmh/api/internal/middleware"
	"dmh/api/internal/service"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	clientIP := l.getClientIP()
	userAgent := l.getUserAgent()

	// 参数验证
	if err := l.validateLoginRequest(req); err != nil {
		// 记录登录失败尝试
		l.svcCtx.AuditService.LogLoginAttempt(nil, req.Username, clientIP, userAgent, false, err.Error())
		return nil, err
	}

	// 查询用户信息
	user, err := l.getUserByUsername(req.Username)
	if err != nil {
		l.Logger.Errorf("查询用户失败: %v", err)
		// 记录登录失败尝试
		l.svcCtx.AuditService.LogLoginAttempt(nil, req.Username, clientIP, userAgent, false, "用户不存在")
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 检查用户状态
	if err := l.checkUserStatus(user); err != nil {
		// 记录登录失败尝试
		l.svcCtx.AuditService.LogLoginAttempt(&user.Id, req.Username, clientIP, userAgent, false, err.Error())
		return nil, err
	}

	// 验证密码
	if err := l.verifyPassword(req.Password, user.Password); err != nil {
		// 记录登录失败
		l.recordLoginAttempt(user.Id, false, clientIP)
		l.svcCtx.AuditService.LogLoginAttempt(&user.Id, req.Username, clientIP, userAgent, false, "密码错误")
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 检查密码是否过期
	expired, err := l.svcCtx.PasswordService.IsPasswordExpired(user.Id)
	if err != nil {
		l.Logger.Errorf("检查密码过期失败: %v", err)
	} else if expired {
		l.svcCtx.AuditService.LogLoginAttempt(&user.Id, req.Username, clientIP, userAgent, false, "密码已过期")
		return nil, fmt.Errorf("密码已过期，请联系管理员重置密码")
	}

	// 查询用户角色
	roles, err := l.getUserRoles(user.Id)
	if err != nil {
		l.Logger.Errorf("查询用户角色失败: %v", err)
		l.svcCtx.AuditService.LogLoginAttempt(&user.Id, req.Username, clientIP, userAgent, false, "获取用户权限失败")
		return nil, fmt.Errorf("获取用户权限失败")
	}

	// 查询品牌管理员的品牌ID
	var brandIDs []int64
	if l.containsRole(roles, "brand_admin") {
		brandIDs, err = l.getUserBrandIDs(user.Id)
		if err != nil {
			l.Logger.Errorf("查询品牌管理员品牌失败: %v", err)
			l.svcCtx.AuditService.LogLoginAttempt(&user.Id, req.Username, clientIP, userAgent, false, "获取品牌权限失败")
			return nil, fmt.Errorf("获取品牌权限失败")
		}
	}

	// 创建用户会话
	session, err := l.svcCtx.SessionService.CreateSession(user.Id, clientIP, userAgent)
	if err != nil {
		l.Logger.Errorf("创建会话失败: %v", err)
		l.svcCtx.AuditService.LogLoginAttempt(&user.Id, req.Username, clientIP, userAgent, false, "创建会话失败")
		return nil, fmt.Errorf("登录失败，请重试")
	}

	// 生成JWT token
	authMiddleware := middleware.NewAuthMiddleware(l.svcCtx.Config.Auth.AccessSecret)
	token, err := authMiddleware.GenerateToken(user.Id, user.Username, roles, brandIDs)
	if err != nil {
		l.Logger.Errorf("生成token失败: %v", err)
		l.svcCtx.AuditService.LogLoginAttempt(&user.Id, req.Username, clientIP, userAgent, false, "生成token失败")
		return nil, fmt.Errorf("登录失败，请重试")
	}

	// 更新最后登录信息
	l.updateLastLogin(user.Id, clientIP)

	// 记录登录成功
	l.recordLoginAttempt(user.Id, true, clientIP)
	l.svcCtx.AuditService.LogLoginAttempt(&user.Id, req.Username, clientIP, userAgent, true, "")

	// 记录审计日志
	l.svcCtx.AuditService.LogUserAction(
		&service.AuditContext{
			UserID:    &user.Id,
			Username:  user.Username,
			ClientIP:  clientIP,
			UserAgent: userAgent,
		},
		"login",
		"user",
		fmt.Sprintf("%d", user.Id),
		map[string]interface{}{
			"session_id": session.ID,
		},
	)

	// 返回登录响应
	return &types.LoginResp{
		Token:    token,
		UserId:   user.Id,
		Username: user.Username,
		Phone:    user.Phone,
		RealName: user.RealName,
		Roles:    roles,
		BrandIds: brandIDs,
	}, nil
}

// validateLoginRequest 验证登录请求
func (l *LoginLogic) validateLoginRequest(req *types.LoginReq) error {
	if strings.TrimSpace(req.Username) == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if strings.TrimSpace(req.Password) == "" {
		return fmt.Errorf("密码不能为空")
	}
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return fmt.Errorf("用户名长度应在3-50个字符之间")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("密码长度不能少于6位")
	}
	return nil
}

// getUserByUsername 根据用户名查询用户
func (l *LoginLogic) getUserByUsername(username string) (*model.User, error) {
	var user model.User
	
	err := l.svcCtx.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}
	
	return &user, nil
}

// checkUserStatus 检查用户状态
func (l *LoginLogic) checkUserStatus(user *model.User) error {
	// 检查用户是否被禁用
	if user.Status == "disabled" {
		return fmt.Errorf("账号已被禁用，请联系管理员")
	}
	
	// 检查用户是否被锁定
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		return fmt.Errorf("账号已被锁定，请稍后再试")
	}
	
	// 获取密码策略
	policy, err := l.svcCtx.PasswordService.GetPasswordPolicy()
	if err != nil {
		l.Logger.Errorf("获取密码策略失败: %v", err)
		// 使用默认值
		policy = &model.PasswordPolicy{MaxLoginAttempts: 5, LockoutDuration: 30}
	}
	
	// 检查登录失败次数
	if user.LoginAttempts >= policy.MaxLoginAttempts {
		// 锁定账号
		lockUntil := time.Now().Add(time.Duration(policy.LockoutDuration) * time.Minute)
		l.lockUser(user.Id, lockUntil)
		
		// 记录安全事件
		l.svcCtx.AuditService.LogSecurityEvent(
			"account_locked",
			"medium",
			&user.Id,
			user.Username,
			l.getClientIP(),
			l.getUserAgent(),
			fmt.Sprintf("用户 %s 因登录失败次数过多被锁定", user.Username),
			map[string]interface{}{
				"login_attempts": user.LoginAttempts,
				"lockout_duration": policy.LockoutDuration,
			},
		)
		
		return fmt.Errorf("登录失败次数过多，账号已被锁定%d分钟", policy.LockoutDuration)
	}
	
	return nil
}

// verifyPassword 验证密码
func (l *LoginLogic) verifyPassword(inputPassword, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}

// getUserRoles 获取用户角色
func (l *LoginLogic) getUserRoles(userID int64) ([]string, error) {
	var roles []model.Role
	
	err := l.svcCtx.DB.Table("roles r").
		Joins("JOIN user_roles ur ON ur.role_id = r.id").
		Where("ur.user_id = ?", userID).
		Find(&roles).Error
	
	if err != nil {
		return nil, err
	}
	
	var roleCodes []string
	for _, role := range roles {
		roleCodes = append(roleCodes, role.Code)
	}
	
	// 如果没有角色，默认为participant
	if len(roleCodes) == 0 {
		roleCodes = []string{"participant"}
	}
	
	return roleCodes, nil
}

// getUserBrandIDs 获取用户管理的品牌ID
func (l *LoginLogic) getUserBrandIDs(userID int64) ([]int64, error) {
	var userBrands []model.UserBrand
	
	err := l.svcCtx.DB.Where("user_id = ?", userID).Find(&userBrands).Error
	if err != nil {
		return nil, err
	}
	
	var brandIDs []int64
	for _, ub := range userBrands {
		brandIDs = append(brandIDs, ub.BrandId)
	}
	
	return brandIDs, nil
}

// containsRole 检查是否包含指定角色
func (l *LoginLogic) containsRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// updateLastLogin 更新最后登录信息
func (l *LoginLogic) updateLastLogin(userID int64, clientIP string) {
	updates := map[string]interface{}{
		"login_attempts": 0,
		"locked_until":   nil,
		"updated_at":     time.Now(),
	}
	
	err := l.svcCtx.DB.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		l.Logger.Errorf("更新最后登录信息失败: %v", err)
	}
}

// lockUser 锁定用户
func (l *LoginLogic) lockUser(userID int64, lockUntil time.Time) {
	err := l.svcCtx.DB.Model(&model.User{}).Where("id = ?", userID).Update("locked_until", lockUntil).Error
	if err != nil {
		l.Logger.Errorf("锁定用户失败: %v", err)
	}
}

// recordLoginAttempt 记录登录尝试
func (l *LoginLogic) recordLoginAttempt(userID int64, success bool, clientIP string) {
	if success {
		// 登录成功，重置失败次数
		l.svcCtx.DB.Model(&model.User{}).Where("id = ?", userID).Update("login_attempts", 0)
	} else {
		// 登录失败，增加失败次数
		l.svcCtx.DB.Model(&model.User{}).Where("id = ?", userID).Update("login_attempts", gorm.Expr("login_attempts + 1"))
	}
	
	// 记录登录日志（可选）
	l.Logger.Infof("用户 %d 登录尝试: 成功=%v, IP=%s", userID, success, clientIP)
}

// getClientIP 获取客户端IP
func (l *LoginLogic) getClientIP() string {
	// 这里应该从HTTP请求中获取真实IP
	// 由于context中可能没有HTTP请求信息，这里返回默认值
	return "unknown"
}

// getUserAgent 获取用户代理
func (l *LoginLogic) getUserAgent() string {
	// 这里应该从HTTP请求中获取用户代理
	// 由于context中可能没有HTTP请求信息，这里返回默认值
	return "unknown"
}
