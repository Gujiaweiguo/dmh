// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package brand

import (
	"context"
	"fmt"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBrandsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBrandsLogic {
	return &GetBrandsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandsLogic) GetBrands() (resp *types.BrandListResp, err error) {
	// 获取用户信息
	userInfo, err := GetUserInfoFromContext(l.ctx)
	if err != nil {
		l.Logger.Errorf("获取用户信息失败: %v", err)
		return nil, fmt.Errorf("用户信息获取失败")
	}

	var brands []model.Brand
	query := l.svcCtx.DB.Model(&model.Brand{})

	// 根据角色过滤品牌
	if userInfo.Role == "brand_admin" {
		// 品牌管理员只能看到自己管理的品牌（已废弃，但保留兼容性）
		var userBrands []model.UserBrand
		if err = l.svcCtx.DB.Where("user_id = ?", userInfo.UserID).Find(&userBrands).Error; err != nil {
			l.Logger.Errorf("查询用户品牌关联失败: %v", err)
			return nil, fmt.Errorf("查询品牌失败")
		}

		if len(userBrands) == 0 {
			return &types.BrandListResp{
				Total:  0,
				Brands: []types.BrandResp{},
			}, nil
		}

		var brandIds []int64
		for _, ub := range userBrands {
			brandIds = append(brandIds, ub.BrandId)
		}
		query = query.Where("id IN ?", brandIds)
	} else if userInfo.Role == "participant" {
		// 普通用户只能看到激活的品牌
		query = query.Where("status = ?", "active")
	}
	// platform_admin 可以看到所有品牌

	if err = query.Find(&brands).Error; err != nil {
		l.Logger.Errorf("查询品牌列表失败: %v", err)
		return nil, fmt.Errorf("查询品牌失败")
	}

	var brandList []types.BrandResp
	for _, brand := range brands {
		brandList = append(brandList, types.BrandResp{
			Id:          brand.Id,
			Name:        brand.Name,
			Logo:        brand.Logo,
			Description: brand.Description,
			Status:      brand.Status,
			CreatedAt:   brand.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   brand.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.BrandListResp{
		Total:  int64(len(brandList)),
		Brands: brandList,
	}, nil
}
