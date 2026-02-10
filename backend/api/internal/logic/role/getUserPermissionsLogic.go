// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserPermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPermissionsLogic {
	return &GetUserPermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserPermissionsLogic) GetUserPermissions(userId int64) (resp *types.UserPermissionResp, err error) {
	var userRoles []model.Role
	var userBrandIds []int64

	l.svcCtx.DB.Table("roles").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userId).
		Find(&userRoles)

	roleCodes := make([]string, 0, len(userRoles))
	for _, r := range userRoles {
		roleCodes = append(roleCodes, r.Code)
	}

	var permissions []model.Permission
	l.svcCtx.DB.Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userId).
		Find(&permissions)

	permCodes := make([]string, 0, len(permissions))
	for _, p := range permissions {
		permCodes = append(permCodes, p.Code)
	}

	l.svcCtx.DB.Table("user_brands").
		Where("user_id = ?", userId).
		Pluck("brand_id", &userBrandIds)

	resp = &types.UserPermissionResp{
		UserId:      userId,
		Roles:       roleCodes,
		Permissions: permCodes,
		BrandIds:    userBrandIds,
	}

	return resp, nil
}
