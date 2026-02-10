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

type DeleteMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMenuLogic) DeleteMenu(menuId int64) (resp *types.CommonResp, err error) {
	tx := l.svcCtx.DB.Begin()

	if err := tx.Where("id = ? OR parent_id = ?", menuId, menuId).Delete(&model.Menu{}).Error; err != nil {
		tx.Rollback()
		l.Errorf("Failed to delete menu: %v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		l.Errorf("Failed to commit delete menu: %v", err)
		return nil, err
	}

	resp = &types.CommonResp{
		Message: "Menu deleted successfully",
	}

	return resp, nil
}
