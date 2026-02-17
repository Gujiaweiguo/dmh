// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"errors"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(userId int64, req *types.AdminUpdateUserReq) (resp *types.UserInfoResp, err error) {
	if userId <= 0 {
		return nil, errors.New("用户ID无效")
	}

	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	if req.RealName != "" {
		user.RealName = req.RealName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != "" {
		user.Status = req.Status
	}

	if err := l.svcCtx.DB.Save(&user).Error; err != nil {
		l.Errorf("更新用户失败: %v", err)
		return nil, errors.New("更新用户失败")
	}

	resp = &types.UserInfoResp{
		Id:        user.Id,
		Username:  user.Username,
		Phone:     user.Phone,
		RealName:  user.RealName,
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}
