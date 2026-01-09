// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"errors"
	"fmt"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req *types.AdminCreateUserReq) (resp *types.UserInfoResp, err error) {
	// 参数验证
	if req.Username == "" || req.Password == "" || req.Phone == "" {
		return nil, fmt.Errorf("用户名、密码和手机号不能为空")
	}

	// 验证角色
	validRoles := map[string]bool{
		"platform_admin": true,
		"participant":    true,
	}
	if !validRoles[req.Role] {
		return nil, fmt.Errorf("无效的用户角色: %s", req.Role)
	}

	// 检查用户名是否已存在
	var existingUser model.User
	err = l.svcCtx.DB.Where("username = ?", req.Username).First(&existingUser).Error
	if err == nil {
		return nil, fmt.Errorf("用户名已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		l.Logger.Errorf("检查用户名失败: %v", err)
		return nil, fmt.Errorf("检查用户名失败")
	}

	// 检查手机号是否已存在
	err = l.svcCtx.DB.Where("phone = ?", req.Phone).First(&existingUser).Error
	if err == nil {
		return nil, fmt.Errorf("手机号已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		l.Logger.Errorf("检查手机号失败: %v", err)
		return nil, fmt.Errorf("检查手机号失败")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		l.Logger.Errorf("密码加密失败: %v", err)
		return nil, fmt.Errorf("密码加密失败")
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Phone:    req.Phone,
		Email:    req.Email,
		RealName: req.RealName,
		Status:   "active",
		Role:     req.Role,
	}

	// 开启事务
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建用户记录
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		l.Logger.Errorf("创建用户失败: %v", err)
		return nil, fmt.Errorf("创建用户失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("提交事务失败: %v", err)
		return nil, fmt.Errorf("创建用户失败")
	}

	l.Logger.Infof("管理员创建用户成功: %s, 角色: %s", user.Username, user.Role)

	// 返回用户信息
	return &types.UserInfoResp{
		Id:        user.Id,
		Username:  user.Username,
		Phone:     user.Phone,
		Email:     user.Email,
		RealName:  user.RealName,
		Status:    user.Status,
		Roles:     []string{user.Role},
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
