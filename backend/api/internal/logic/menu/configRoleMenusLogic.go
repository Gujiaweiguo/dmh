// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigRoleMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigRoleMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigRoleMenusLogic {
	return &ConfigRoleMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigRoleMenusLogic) ConfigRoleMenus(req *types.RoleMenuReq) (resp *types.CommonResp, err error) {
	tx := l.svcCtx.DB.Begin()

	if err := tx.Where("role_id = ?", req.RoleId).Delete(&model.RoleMenu{}).Error; err != nil {
		tx.Rollback()
		l.Errorf("Failed to delete old role menus: %v", err)
		return nil, err
	}

	for _, menuId := range req.MenuIds {
		roleMenu := &model.RoleMenu{
			RoleID: req.RoleId,
			MenuID: menuId,
		}
		if err := tx.Create(roleMenu).Error; err != nil {
			tx.Rollback()
			l.Errorf("Failed to create role menu: %v", err)
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		l.Errorf("Failed to commit role menus: %v", err)
		return nil, err
	}

	l.Infof("Role menus configured successfully: roleId=%d, menus=%d", req.RoleId, len(req.MenuIds))

	resp = &types.CommonResp{
		Message: "Role menus configured successfully",
	}

	return resp, nil
}
