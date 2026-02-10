// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"errors"

	"dmh/api/internal/middleware"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"
	"golang.org/x/crypto/bcrypt"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordReq) (resp *types.CommonResp, err error) {
	// 从context获取userId
	userId, err := middleware.GetUserIDFromContext(l.ctx)
	if err != nil {
		return nil, err
	}
	if userId == 0 {
		return nil, errors.New("未登录")
	}

	// 查询用户
	var user model.User
	err = l.svcCtx.DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		l.Errorf("查询用户失败: %v", err)
		return nil, errors.New("用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return nil, errors.New("旧密码错误")
	}

	// 简单的密码强度检查
	if len(req.NewPassword) < 6 {
		return nil, errors.New("新密码长度不能少于6位")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		l.Errorf("密码加密失败: %v", err)
		return nil, errors.New("密码加密失败")
	}

	// 更新密码
	err = l.svcCtx.DB.Model(&user).Update("password", string(hashedPassword)).Error
	if err != nil {
		l.Errorf("更新密码失败: %v", err)
		return nil, errors.New("更新密码失败")
	}

	resp = &types.CommonResp{
		Message: "密码修改成功",
	}

	return resp, nil
}
