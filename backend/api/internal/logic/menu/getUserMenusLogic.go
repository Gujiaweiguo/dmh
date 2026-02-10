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

type GetUserMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserMenusLogic {
	return &GetUserMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserMenusLogic) GetUserMenus(userId int64, platform string) (resp *types.UserMenuResp, err error) {
	var roleMenus []model.RoleMenu

	l.svcCtx.DB.Table("role_menus").
		Joins("JOIN user_roles ON user_roles.role_id = role_menus.role_id").
		Where("user_roles.user_id = ?", userId).
		Find(&roleMenus)

	menuIds := make([]int64, 0, len(roleMenus))
	for _, rm := range roleMenus {
		menuIds = append(menuIds, rm.MenuID)
	}

	var menus []model.Menu
	query := l.svcCtx.DB.Where("platform = ? AND status = ?", platform, "active")
	if len(menuIds) > 0 {
		query = query.Where("id IN ?", menuIds)
	}

	if err := query.Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		l.Errorf("Failed to get user menus: %v", err)
		return nil, err
	}

	menuMap := make(map[int64]*types.MenuResp)
	var rootMenus []types.MenuResp

	for _, m := range menus {
		menuResp := &types.MenuResp{
			Id:        m.ID,
			Name:      m.Name,
			Code:      m.Code,
			Path:      m.Path,
			Icon:      m.Icon,
			Sort:      m.Sort,
			Type:      m.Type,
			Platform:  m.Platform,
			Status:    m.Status,
			CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if m.ParentID != nil {
			menuResp.ParentId = *m.ParentID
		}

		menuMap[m.ID] = menuResp
	}

	for _, m := range menus {
		if m.ParentID == nil {
			rootMenus = append(rootMenus, *menuMap[m.ID])
		} else if parent, exists := menuMap[*m.ParentID]; exists {
			parent.Children = append(parent.Children, *menuMap[m.ID])
		}
	}

	resp = &types.UserMenuResp{
		UserId:   userId,
		Platform: platform,
		Menus:    rootMenus,
	}

	return resp, nil
}
