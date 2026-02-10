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

type GetRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRolesLogic {
	return &GetRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRolesLogic) GetRoles() (resp []types.RoleResp, err error) {
	var roles []model.Role

	if err := l.svcCtx.DB.Find(&roles).Error; err != nil {
		l.Errorf("Failed to get roles: %v", err)
		return nil, err
	}

	resp = make([]types.RoleResp, 0, len(roles))
	for _, r := range roles {
		roleResp := types.RoleResp{
			Id:          r.ID,
			Name:        r.Name,
			Code:        r.Code,
			Description: r.Description,
			CreatedAt:   r.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		var permissions []model.Permission
		l.svcCtx.DB.Table("permissions").
			Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
			Where("role_permissions.role_id = ?", r.ID).
			Find(&permissions)

		permCodes := make([]string, 0, len(permissions))
		for _, p := range permissions {
			permCodes = append(permCodes, p.Code)
		}
		roleResp.Permissions = permCodes

		resp = append(resp, roleResp)
	}

	return resp, nil
}
