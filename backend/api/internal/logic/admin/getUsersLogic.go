// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersLogic {
	return &GetUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUsersLogic) GetUsers(req *types.AdminGetUsersReq) (resp *types.AdminUserListResp, err error) {
	var users []model.User
	var total int64

	// 构建查询
	query := l.svcCtx.DB.Model(&model.User{})

	// 添加筛选条件
	if req.Role != "" {
		query = query.Where("role = ?", req.Role)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Keyword != "" {
		query = query.Where("username LIKE ? OR phone LIKE ? OR real_name LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if req.Page > 0 && req.PageSize > 0 {
		query = query.Offset(int(offset)).Limit(int(req.PageSize))
	}

	// 执行查询
	err = query.Find(&users).Error
	if err != nil {
		l.Errorf("查询用户列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	userInfos := make([]types.UserInfoResp, 0, len(users))
	for _, user := range users {
		// 获取用户角色
		var roles []model.Role
		l.svcCtx.DB.Table("roles").
			Joins("JOIN user_roles ON user_roles.role_id = roles.id").
			Where("user_roles.user_id = ?", user.Id).
			Find(&roles)

		roleCodes := make([]string, 0, len(roles))
		for _, role := range roles {
			roleCodes = append(roleCodes, role.Code)
		}

		userInfos = append(userInfos, types.UserInfoResp{
			Id:        user.Id,
			Username:  user.Username,
			Phone:     user.Phone,
			Email:     user.Email,
			RealName:  user.RealName,
			Status:    user.Status,
			Roles:     roleCodes,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	resp = &types.AdminUserListResp{
		Total: total,
		Users: userInfos,
	}

	return resp, nil
}
