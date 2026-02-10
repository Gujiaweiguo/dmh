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

type GetMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuLogic {
	return &GetMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenuLogic) GetMenu(menuId int64) (resp *types.MenuResp, err error) {
	menu := &model.Menu{}
	if err := l.svcCtx.DB.First(menu, menuId).Error; err != nil {
		l.Errorf("Menu not found: %v", err)
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
