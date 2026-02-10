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

type CreateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMenuLogic) CreateMenu(req *types.CreateMenuReq) (resp *types.MenuResp, err error) {
	menu := &model.Menu{
		Name:     req.Name,
		Code:     req.Code,
		Path:     req.Path,
		Icon:     req.Icon,
		ParentID: nil,
		Sort:     req.Sort,
		Type:     req.Type,
		Platform: req.Platform,
		Status:   "active",
	}

	if req.ParentId > 0 {
		menu.ParentID = &req.ParentId
	}

	if err := l.svcCtx.DB.Create(menu).Error; err != nil {
		l.Errorf("Failed to create menu: %v", err)
		return nil, err
	}

	resp = &types.MenuResp{
		Id:        menu.ID,
		Name:      menu.Name,
		Code:      menu.Code,
		Path:      menu.Path,
		Icon:      menu.Icon,
		ParentId:  req.ParentId,
		Sort:      menu.Sort,
		Type:      menu.Type,
		Platform:  menu.Platform,
		Status:    menu.Status,
		CreatedAt: menu.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}
