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

type UpdateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMenuLogic) UpdateMenu(menuId int64, req *types.UpdateMenuReq) (resp *types.MenuResp, err error) {
	menu := &model.Menu{}
	if err := l.svcCtx.DB.First(menu, menuId).Error; err != nil {
		l.Errorf("Menu not found: %v", err)
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Code != "" {
		updates["code"] = req.Code
	}
	if req.Path != "" {
		updates["path"] = req.Path
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.Sort > 0 {
		updates["sort"] = req.Sort
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Platform != "" {
		updates["platform"] = req.Platform
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.ParentId > 0 {
		updates["parent_id"] = req.ParentId
	}

	if err := l.svcCtx.DB.Model(menu).Updates(updates).Error; err != nil {
		l.Errorf("Failed to update menu: %v", err)
		return nil, err
	}

	resp = &types.MenuResp{
		Id:        menu.ID,
		Name:      menu.Name,
		Code:      menu.Code,
		Path:      menu.Path,
		Icon:      menu.Icon,
		Sort:      menu.Sort,
		Type:      menu.Type,
		Platform:  menu.Platform,
		Status:    menu.Status,
		CreatedAt: menu.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if menu.ParentID != nil {
		resp.ParentId = *menu.ParentID
	}

	return resp, nil
}
