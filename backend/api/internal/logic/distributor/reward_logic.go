package distributor

import (
	"context"
	"fmt"
	"time"

	"dmh/api/internal/svc"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RewardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRewardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RewardLogic {
	return &RewardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DistributeReward 分销商奖励分配
// 在订单支付成功后调用，向上追溯最多3级分销商并分配奖励
func (l *RewardLogic) DistributeReward(orderId int64, amount float64, brandId int64) error {
	// 幂等性检查：检查该订单是否已经分配过奖励
	var existingRewardCount int64
	l.svcCtx.DB.Model(&model.DistributorReward{}).Where("order_id = ?", orderId).Count(&existingRewardCount)
	if existingRewardCount > 0 {
		l.Logger.Infof("订单 %d 已分配过奖励，跳过", orderId)
		return nil
	}

	// 获取订单信息
	var order model.Order
	if err := l.svcCtx.DB.Where("id = ?", orderId).First(&order).Error; err != nil {
		return fmt.Errorf("订单不存在")
	}

	// 如果订单有推荐人，从推荐人开始追溯
	// 否则跳过
	if order.ReferrerId == 0 {
		return nil // 没有推荐人，无需分配奖励
	}

	// 开始事务
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 事务内再次幂等性检查（双重检查）
	var count int64
	tx.Model(&model.DistributorReward{}).Where("order_id = ?", orderId).Count(&count)
	if count > 0 {
		tx.Rollback()
		l.Logger.Infof("订单 %d 已分配过奖励（事务内检查），跳过", orderId)
		return nil
	}

	// 获取品牌奖励配置
	levelRewards, err := l.getBrandLevelRewards(tx, brandId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("获取奖励配置失败")
	}

	// 向上追溯分销商，最多3级
	currentUserId := order.ReferrerId
	level := 1

	for level <= 3 && currentUserId > 0 {
		// 查找当前用户是否是分销商
		var distributor model.Distributor
		err := tx.Where("user_id = ? AND brand_id = ? AND status = ?", currentUserId, brandId, "active").
			First(&distributor).Error

		if err != nil {
			// 该用户不是分销商，继续向上查找
			nextParentId := l.getParentDistributorId(tx, currentUserId, brandId)
			if nextParentId == 0 {
				break // 没有上级分销商了
			}
			currentUserId = nextParentId
			continue
		}

		// 获取该级别的奖励比例
		rewardRate, ok := levelRewards[level]
		if !ok || rewardRate <= 0 {
			// 该级别没有配置奖励比例，跳过
			currentUserId = l.getParentDistributorId(tx, distributor.Id, brandId)
			if currentUserId == 0 {
				break
			}
			// 注意：这里需要获取的是上级分销商的user_id，而不是distributor_id
			var parentDistributor model.Distributor
			if distributor.ParentId != nil {
				tx.Where("id = ?", *distributor.ParentId).First(&parentDistributor)
				currentUserId = parentDistributor.UserId
			} else {
				currentUserId = 0
			}
			continue
		}

		// 计算奖励金额
		rewardAmount := amount * rewardRate / 100

		// 创建奖励记录
		now := time.Now()
		reward := model.DistributorReward{
			DistributorId: distributor.Id,
			UserId:        distributor.UserId,
			OrderId:       order.Id,
			CampaignId:    order.CampaignId,
			Amount:        rewardAmount,
			Level:         level,
			RewardRate:    rewardRate,
			FromUserId:    order.MemberID, // 购买用户的会员ID
			Status:        "settled",
			SettledAt:     &now,
		}

		if err := tx.Create(&reward).Error; err != nil {
			tx.Rollback()
			l.Logger.Errorf("创建奖励记录失败: %v", err)
			return fmt.Errorf("创建奖励记录失败")
		}

		// 更新分销商累计收益
		if err := tx.Model(&distributor).
			UpdateColumn("total_earnings", gorm.Expr("total_earnings + ?", rewardAmount)).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("更新分销商收益失败")
		}

		// 更新用户余额
		if err := l.updateUserBalance(tx, distributor.UserId, rewardAmount); err != nil {
			tx.Rollback()
			return err
		}

		// 继续查找上级分销商
		if distributor.ParentId != nil {
			var parentDistributor model.Distributor
			if err := tx.Where("id = ?", *distributor.ParentId).First(&parentDistributor).Error; err == nil {
				currentUserId = parentDistributor.UserId
			} else {
				break
			}
		} else {
			break
		}

		level++
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败")
	}

	return nil
}

// DistributeRewardByCode 根据推广码分配奖励
// 当用户通过分销商推广链接下单时调用
func (l *RewardLogic) DistributeRewardByCode(linkCode string, orderId int64, amount float64) error {
	// 获取分销商信息
	var link model.DistributorLink
	if err := l.svcCtx.DB.Where("link_code = ?", linkCode).First(&link).Error; err != nil {
		return fmt.Errorf("推广链接不存在")
	}

	// 获取分销商
	var distributor model.Distributor
	if err := l.svcCtx.DB.Where("id = ? AND status = ?", link.DistributorId, "active").First(&distributor).Error; err != nil {
		return fmt.Errorf("分销商不存在或未激活")
	}

	// 获取订单信息
	var order model.Order
	if err := l.svcCtx.DB.Where("id = ?", orderId).First(&order).Error; err != nil {
		return fmt.Errorf("订单不存在")
	}

	// 更新订单的推荐人ID
	if order.ReferrerId == 0 {
		l.svcCtx.DB.Model(&order).Update("referrer_id", distributor.UserId)
	}

	// 使用 DistributeReward 进行多级奖励分配
	return l.DistributeReward(orderId, amount, distributor.BrandId)
}

// getBrandLevelRewards 获取品牌各级奖励配置
func (l *RewardLogic) getBrandLevelRewards(db *gorm.DB, brandId int64) (map[int]float64, error) {
	var rewards []model.DistributorLevelReward
	if err := db.Where("brand_id = ?", brandId).Find(&rewards).Error; err != nil {
		return nil, err
	}

	result := make(map[int]float64)
	for _, r := range rewards {
		result[r.Level] = r.RewardPercentage
	}

	// 如果没有配置，使用默认值
	if len(result) == 0 {
		result[1] = 5.0 // 一级5%
		result[2] = 2.0 // 二级2%
		result[3] = 1.0 // 三级1%
	}

	return result, nil
}

// getParentDistributorId 获取上级分销商的user_id
func (l *RewardLogic) getParentDistributorId(db *gorm.DB, userId, brandId int64) int64 {
	var distributor model.Distributor
	if err := db.Where("user_id = ? AND brand_id = ?", userId, brandId).First(&distributor).Error; err != nil {
		return 0
	}

	if distributor.ParentId != nil {
		var parent model.Distributor
		if err := db.Where("id = ?", *distributor.ParentId).First(&parent).Error; err == nil {
			return parent.UserId
		}
	}

	return 0
}

// updateUserBalance 更新用户余额
func (l *RewardLogic) updateUserBalance(db *gorm.DB, userId int64, amount float64) error {
	var balance model.UserBalance
	err := db.Where("user_id = ?", userId).First(&balance).Error

	if err == gorm.ErrRecordNotFound {
		// 创建余额记录
		balance = model.UserBalance{
			UserId:      userId,
			Balance:     amount,
			TotalReward: amount,
			Version:     0,
		}
		return db.Create(&balance).Error
	}

	if err != nil {
		return err
	}

	// 使用乐观锁更新余额
	return db.Model(&balance).
		Where("version = ?", balance.Version).
		Updates(map[string]interface{}{
			"balance":      gorm.Expr("balance + ?", amount),
			"total_reward": gorm.Expr("total_reward + ?", amount),
			"version":      gorm.Expr("version + 1"),
		}).Error
}

// LinkUserToDistributor 将用户关联到分销商（通过推广链接访问时）
func (l *RewardLogic) LinkUserToDistributor(linkCode string, userId int64) error {
	// 查找推广链接
	var link model.DistributorLink
	if err := l.svcCtx.DB.Where("link_code = ?", linkCode).First(&link).Error; err != nil {
		return fmt.Errorf("推广链接不存在")
	}

	// 查找分销商
	var distributor model.Distributor
	if err := l.svcCtx.DB.Where("id = ?", link.DistributorId).First(&distributor).Error; err != nil {
		return fmt.Errorf("分销商不存在")
	}

	// 存储用户与分销商的关联到session或cookie中
	// 这里简单返回分销商信息，由调用方处理
	// 实际应用中可以使用临时表或Redis存储
	return nil
}

// GetDistributorChain 获取用户到分销商的完整链条
func (l *RewardLogic) GetDistributorChain(userId, brandId int64) ([]model.Distributor, error) {
	chain := make([]model.Distributor, 0)
	maxLevel := 3

	currentUserId := userId
	for i := 0; i < maxLevel; i++ {
		var distributor model.Distributor
		err := l.svcCtx.DB.Preload("User").Preload("Brand").
			Where("user_id = ? AND brand_id = ?", currentUserId, brandId).
			First(&distributor).Error

		if err != nil {
			// 不是分销商，链条结束
			break
		}

		chain = append(chain, distributor)

		if distributor.ParentId != nil {
			var parent model.Distributor
			if err := l.svcCtx.DB.Where("id = ?", *distributor.ParentId).First(&parent).Error; err != nil {
				break
			}
			currentUserId = parent.UserId
		} else {
			break
		}
	}

	return chain, nil
}
