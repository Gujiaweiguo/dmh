package distributor

import (
	"context"
	"dmh/api/internal/svc"
	"dmh/model"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// AutoUpgradeLogic 自动成为分销商业务逻辑
type AutoUpgradeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewAutoUpgradeLogic 创建自动升级逻辑实例
func NewAutoUpgradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AutoUpgradeLogic {
	return &AutoUpgradeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CheckAndAutoUpgrade 检查并自动成为分销商
// 在订单支付成功后调用，检查用户是否已成为分销商，如未成为则自动创建
func (l *AutoUpgradeLogic) CheckAndAutoUpgrade(userID int64, brandID int64, referrerID int64) (*model.Distributor, error) {
	// 1. 检查用户是否已是该品牌的分销商
	var existingDistributor model.Distributor
	err := l.svcCtx.DB.Where("user_id = ? AND brand_id = ?", userID, brandID).First(&existingDistributor).Error
	if err == nil {
		// 已是分销商，直接返回
		return &existingDistributor, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 数据库错误
		return nil, fmt.Errorf("查询分销商信息失败: %w", err)
	}

	// 2. 查询活动是否启用分销
	// 注意：这里需要传入campaignID，暂时假设所有活动都支持分销
	// 实际使用时需要根据订单的campaign_id查询活动的enable_distribution字段

	// 3. 检查推荐人是否是分销商
	var parentID int64 = 0
	if referrerID > 0 {
		var parentDistributor model.Distributor
		err := l.svcCtx.DB.Where("user_id = ? AND brand_id = ? AND status = ?", referrerID, brandID, "active").
			First(&parentDistributor).Error
		if err == nil {
			// 推荐人是分销商，设置为上级
			parentID = parentDistributor.Id
		}
	}

	// 4. 创建分销商记录
	distributor := &model.Distributor{
		UserId:            userID,
		BrandId:           brandID,
		Level:             1,         // 默认一级
		ParentId:          &parentID, // 上级分销商ID
		Status:            "active",  // 默认激活状态
		TotalEarnings:     0,
		SubordinatesCount: 0,
	}

	if err := l.svcCtx.DB.Create(distributor).Error; err != nil {
		return nil, fmt.Errorf("创建分销商记录失败: %w", err)
	}

	// 5. 添加distributor角色
	if err := l.addDistributorRole(userID); err != nil {
		// 添加角色失败，但分销商记录已创建，记录日志
		return distributor, fmt.Errorf("分销商创建成功，但添加角色失败: %w", err)
	}

	// 6. 更新上级分销商的下级数量（如果有上级）
	if parentID > 0 {
		l.updateSubordinatesCount(parentID)
	}

	return distributor, nil
}

// CheckAndAutoUpgradeWithCampaign 检查并自动成为分销商（带活动ID）
func (l *AutoUpgradeLogic) CheckAndAutoUpgradeWithCampaign(userID int64, brandID int64, campaignID int64, referrerID int64) (*model.Distributor, error) {
	// 1. 检查用户是否已是该品牌的分销商
	var existingDistributor model.Distributor
	err := l.svcCtx.DB.Where("user_id = ? AND brand_id = ?", userID, brandID).First(&existingDistributor).Error
	if err == nil {
		// 已是分销商，直接返回
		return &existingDistributor, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 数据库错误
		return nil, fmt.Errorf("查询分销商信息失败: %w", err)
	}

	// 2. 查询活动是否启用分销
	var campaign model.Campaign
	err = l.svcCtx.DB.Where("id = ? AND enable_distribution = ?", campaignID, true).First(&campaign).Error
	if err != nil {
		// 活动未启用分销或活动不存在，不创建分销商
		return nil, fmt.Errorf("活动未启用分销")
	}

	// 3. 检查推荐人是否是分销商
	var parentID int64 = 0
	if referrerID > 0 {
		var parentDistributor model.Distributor
		err := l.svcCtx.DB.Where("user_id = ? AND brand_id = ? AND status = ?", referrerID, brandID, "active").
			First(&parentDistributor).Error
		if err == nil {
			// 推荐人是分销商，设置为上级
			parentID = parentDistributor.Id
		}
	}

	// 4. 创建分销商记录
	distributor := &model.Distributor{
		UserId:            userID,
		BrandId:           brandID,
		Level:             1,         // 默认一级
		ParentId:          &parentID, // 上级分销商ID
		Status:            "active",  // 默认激活状态
		TotalEarnings:     0,
		SubordinatesCount: 0,
	}

	if err := l.svcCtx.DB.Create(distributor).Error; err != nil {
		return nil, fmt.Errorf("创建分销商记录失败: %w", err)
	}

	// 5. 添加distributor角色
	if err := l.addDistributorRole(userID); err != nil {
		// 添加角色失败，但分销商记录已创建，记录日志
		return distributor, fmt.Errorf("分销商创建成功，但添加角色失败: %w", err)
	}

	// 6. 更新上级分销商的下级数量（如果有上级）
	if parentID > 0 {
		l.updateSubordinatesCount(parentID)
	}

	return distributor, nil
}

// addDistributorRole 为用户添加distributor角色
func (l *AutoUpgradeLogic) addDistributorRole(userID int64) error {
	// 查询distributor角色
	var role model.Role
	err := l.svcCtx.DB.Where("code = ?", "distributor").First(&role).Error
	if err != nil {
		return fmt.Errorf("查询distributor角色失败: %w", err)
	}

	// 检查用户是否已有该角色
	var userRole model.UserRole
	err = l.svcCtx.DB.Where("user_id = ? AND role_id = ?", userID, role.ID).First(&userRole).Error
	if err == nil {
		// 已有该角色，无需重复添加
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("查询用户角色失败: %w", err)
	}

	// 添加角色
	userRole = model.UserRole{
		UserID: userID,
		RoleID: role.ID,
	}
	if err := l.svcCtx.DB.Create(&userRole).Error; err != nil {
		return fmt.Errorf("添加用户角色失败: %w", err)
	}

	return nil
}

// updateSubordinatesCount 更新分销商的下级数量
func (l *AutoUpgradeLogic) updateSubordinatesCount(distributorID int64) error {
	// 查询下级数量
	var count int64
	err := l.svcCtx.DB.Model(&model.Distributor{}).Where("parent_id = ?", distributorID).Count(&count).Error
	if err != nil {
		return fmt.Errorf("查询下级数量失败: %w", err)
	}

	// 更新下级数量
	err = l.svcCtx.DB.Model(&model.Distributor{}).Where("id = ?", distributorID).
		Update("subordinates_count", count).Error
	if err != nil {
		return fmt.Errorf("更新下级数量失败: %w", err)
	}

	return nil
}

// CalculateDistributorPath 计算分销链路径
// 返回格式：一级ID,二级ID,三级ID
func (l *AutoUpgradeLogic) CalculateDistributorPath(distributorID int64, maxLevel int) string {
	if distributorID == 0 {
		return ""
	}

	path := []int64{}
	currentID := distributorID
	level := 0

	for currentID > 0 && level < maxLevel {
		var distributor model.Distributor
		err := l.svcCtx.DB.Where("id = ?", currentID).First(&distributor).Error
		if err != nil {
			break
		}

		path = append(path, distributor.Id)
		level++

		if distributor.ParentId == nil {
			break
		}
		currentID = *distributor.ParentId
	}

	// 反转路径（从上级到下级）
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	// 转换为字符串
	return fmt.Sprintf("%v", path)
}
