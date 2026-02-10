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

type GetPermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPermissionsLogic {
	return &GetPermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPermissionsLogic) GetPermissions() (resp []types.PermissionResp, err error) {
	var permissions []model.Permission

	if err := l.svcCtx.DB.Find(&permissions).Error; err != nil {
		l.Errorf("Failed to get permissions: %v", err)
		return nil, err
	}

	resp = make([]types.PermissionResp, 0, len(permissions))
	for _, p := range permissions {
		resp = append(resp, types.PermissionResp{
			Id:          p.ID,
			Name:        p.Name,
			Code:        p.Code,
			Resource:    p.Resource,
			Action:      p.Action,
			Description: p.Description,
		})
	}

	return resp, nil
}
