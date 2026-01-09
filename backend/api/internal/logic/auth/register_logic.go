// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"dmh/api/internal/middleware"
	"dmh/api/internal/service"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.LoginResp, err error) {
	clientIP := l.getClientIP()
	userAgent := l.getUserAgent()

	// 参数验证
	if err := l.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// 使用密码服务验证密码强度
	if err := l.svcCtx.PasswordService.ValidatePassword(req.Password, 0); err != nil {
		return nil, err
	}

	// 检查用户名和手机号唯一性
	if err := l.checkUniqueness(req.Username, req.Phone); err != nil {
		return nil, err
	}

	// 使用密码服务加密密码
	hashedPassword, err := l.svcCtx.PasswordService.HashPassword(req.Password)
	if err != nil {
		l.Logger.Errorf("密码加密失败: %v", err)
		return nil, fmt.Errorf("注册失败，请重试")
	}

	// 开启数据库事务
	tx := l.svcCtx.DB.Begin()
	if tx.Error != nil {
		l.Logger.Errorf("开启事务失败: %v", tx.Error)
		return nil, fmt.Errorf("注册失败，请重试")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建用户记录
	userID, err := l.createUser(tx, req, hashedPassword)
	if err != nil {
		tx.Rollback()
		l.Logger.Errorf("创建用户失败: %v", err)
		return nil, fmt.Errorf("注册失败，请重试")
	}

	// 分配默认角色 (participant)
	if err := l.assignDefaultRole(tx, userID); err != nil {
		tx.Rollback()
		l.Logger.Errorf("分配默认角色失败: %v", err)
		return nil, fmt.Errorf("注册失败，请重试")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("提交事务失败: %v", err)
		return nil, fmt.Errorf("注册失败，请重试")
	}

	// 保存密码历史
	if err := l.svcCtx.PasswordService.SavePasswordHistory(userID, hashedPassword); err != nil {
		l.Logger.Errorf("保存密码历史失败: %v", err)
		// 不影响注册流程，继续执行
	}

	// 生成JWT token
	authMiddleware := middleware.NewAuthMiddleware(l.svcCtx.Config.Auth.AccessSecret)
	roles := []string{"participant"}
	token, err := authMiddleware.GenerateToken(userID, req.Username, roles, []int64{})
	if err != nil {
		l.Logger.Errorf("生成token失败: %v", err)
		return nil, fmt.Errorf("注册成功，但登录失败，请手动登录")
	}

	// 记录审计日志
	l.svcCtx.AuditService.LogUserAction(
		&service.AuditContext{
			UserID:    &userID,
			Username:  req.Username,
			ClientIP:  clientIP,
			UserAgent: userAgent,
		},
		"register",
		"user",
		fmt.Sprintf("%d", userID),
		map[string]interface{}{
			"phone": req.Phone,
			"email": req.Email,
		},
	)

	l.Logger.Infof("用户注册成功: ID=%d, Username=%s", userID, req.Username)

	// 返回注册响应
	return &types.LoginResp{
		Token:    token,
		UserId:   userID,
		Username: req.Username,
		Phone:    req.Phone,
		RealName: req.RealName,
		Roles:    roles,
		BrandIds: []int64{},
	}, nil
}

// validateRegisterRequest 验证注册请求
func (l *RegisterLogic) validateRegisterRequest(req *types.RegisterReq) error {
	// 用户名验证
	if strings.TrimSpace(req.Username) == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return fmt.Errorf("用户名长度应在3-50个字符之间")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(req.Username) {
		return fmt.Errorf("用户名只能包含字母、数字和下划线")
	}

	// 密码验证
	if strings.TrimSpace(req.Password) == "" {
		return fmt.Errorf("密码不能为空")
	}
	if len(req.Password) < 6 || len(req.Password) > 50 {
		return fmt.Errorf("密码长度应在6-50个字符之间")
	}
	if !l.isValidPassword(req.Password) {
		return fmt.Errorf("密码必须包含字母和数字")
	}

	// 手机号验证
	if strings.TrimSpace(req.Phone) == "" {
		return fmt.Errorf("手机号不能为空")
	}
	if !l.isValidPhone(req.Phone) {
		return fmt.Errorf("手机号格式不正确")
	}

	// 邮箱验证（可选）
	if req.Email != "" && !l.isValidEmail(req.Email) {
		return fmt.Errorf("邮箱格式不正确")
	}

	// 真实姓名验证（可选）
	if req.RealName != "" && len(req.RealName) > 50 {
		return fmt.Errorf("真实姓名长度不能超过50个字符")
	}

	return nil
}

// isValidPassword 验证密码强度
func (l *RegisterLogic) isValidPassword(password string) bool {
	// 至少包含一个字母和一个数字
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	return hasLetter && hasNumber
}

// isValidPhone 验证手机号格式
func (l *RegisterLogic) isValidPhone(phone string) bool {
	// 中国大陆手机号格式
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// isValidEmail 验证邮箱格式
func (l *RegisterLogic) isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// checkUniqueness 检查用户名和手机号唯一性
func (l *RegisterLogic) checkUniqueness(username, phone string) error {
	// 检查用户名是否已存在
	var count int64
	err := l.svcCtx.DB.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return fmt.Errorf("检查用户名唯一性失败")
	}
	if count > 0 {
		return fmt.Errorf("用户名已存在")
	}

	// 检查手机号是否已存在
	err = l.svcCtx.DB.Model(&model.User{}).Where("phone = ?", phone).Count(&count).Error
	if err != nil {
		return fmt.Errorf("检查手机号唯一性失败")
	}
	if count > 0 {
		return fmt.Errorf("手机号已被注册")
	}

	return nil
}

// hashPassword 加密密码
func (l *RegisterLogic) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// createUser 创建用户记录
func (l *RegisterLogic) createUser(tx *gorm.DB, req *types.RegisterReq, hashedPassword string) (int64, error) {
	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Phone:    req.Phone,
		Email:    req.Email,
		RealName: req.RealName,
		Status:   "active",
	}
	
	err := tx.Create(user).Error
	if err != nil {
		return 0, err
	}
	
	return user.Id, nil
}

// assignDefaultRole 分配默认角色
func (l *RegisterLogic) assignDefaultRole(tx *gorm.DB, userID int64) error {
	// 查询participant角色ID
	var role model.Role
	err := tx.Where("code = ?", "participant").First(&role).Error
	if err != nil {
		return fmt.Errorf("查询participant角色失败: %v", err)
	}
	
	// 分配角色
	userRole := &model.UserRole{
		UserID: userID,
		RoleID: role.ID,
	}
	
	err = tx.Create(userRole).Error
	if err != nil {
		return fmt.Errorf("分配角色失败: %v", err)
	}
	
	return nil
}

// getClientIP 获取客户端IP
func (l *RegisterLogic) getClientIP() string {
	// 这里应该从HTTP请求中获取真实IP
	// 由于context中可能没有HTTP请求信息，这里返回默认值
	return "unknown"
}

// getUserAgent 获取用户代理
func (l *RegisterLogic) getUserAgent() string {
	// 这里应该从HTTP请求中获取用户代理
	// 由于context中可能没有HTTP请求信息，这里返回默认值
	return "unknown"
}