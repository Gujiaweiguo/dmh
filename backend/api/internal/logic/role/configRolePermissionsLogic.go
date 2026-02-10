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

type ConfigRolePermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigRolePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigRolePermissionsLogic {
	return &ConfigRolePermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigRolePermissionsLogic) ConfigRolePermissions(req *types.RolePermissionReq) (resp *types.CommonResp, err error) {
	tx := l.svcCtx.DB.Begin()

	if err := tx.Where("role_id = ?", req.RoleId).Delete(&model.RolePermission{}).Error; err != nil {
		tx.Rollback()
		l.Errorf("Failed to delete old role permissions: %v", err)
		return nil, err
	}

	for _, permId := range req.PermissionIds {
		rolePerm := &model.RolePermission{
			RoleID:       req.RoleId,
			PermissionID: permId,
		}
		if err := tx.Create(rolePerm).Error; err != nil {
			tx.Rollback()
			l.Errorf("Failed to create role permission: %v", err)
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		l.Errorf("Failed to commit role permissions: %v", err)
		return nil, err
	}

	l.Infof("Role permissions configured successfully: roleId=%d, permissions=%d", req.RoleId, len(req.PermissionIds))

	resp = &types.CommonResp{
		Message: "Role permissions configured successfully",
	}

	return resp, nil
}
