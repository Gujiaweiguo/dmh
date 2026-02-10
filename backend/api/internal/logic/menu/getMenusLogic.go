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

type GetMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenusLogic {
	return &GetMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenusLogic) GetMenus(req *types.GetMenusReq) (resp *types.MenuListResp, err error) {
	var menus []model.Menu
	var total int64

	query := l.svcCtx.DB.Model(&model.Menu{})

	if req.Platform != "" {
		query = query.Where("platform = ?", req.Platform)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	if err := query.Count(&total).Error; err != nil {
		l.Errorf("Failed to count menus: %v", err)
		return nil, err
	}

	if err := query.Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		l.Errorf("Failed to get menus: %v", err)
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

	resp = &types.MenuListResp{
		Total: total,
		Menus: rootMenus,
	}

	return resp, nil
}
