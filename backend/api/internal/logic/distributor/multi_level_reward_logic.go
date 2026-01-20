package distributor

import (
	"context"
	"dmh/api/internal/svc"
	"dmh/model"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// MultiLevelRewardLogic 多级奖励计算业务逻辑
type MultiLevelRewardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewMultiLevelRewardLogic 创建多级奖励逻辑实例
func NewMultiLevelRewardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiLevelRewardLogic {
	return &MultiLevelRewardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DistributionRewardConfig 分销奖励配置
type DistributionRewardConfig struct {
	Level1 float64 `json:"level1"` // 一级奖励比例
	Level2 float64 `json:"level2"` // 二级奖励比例
	Level3 float64 `json:"level3"` // 三级奖励比例
}

// CalculateAndDistributeRewards 计算并分配多级分销奖励
// 在订单支付成功后调用
func (l *MultiLevelRewardLogic) CalculateAndDistributeRewards(orderID int64, campaignID int64, userID int64, referrerID int64, orderAmount float64) error {
	// 幂等性检查：检查该订单是否已经分配过奖励
	var existingRewardCount int64
	l.svcCtx.DB.Model(&model.DistributorReward{}).Where("order_id = ?", orderID).Count(&existingRewardCount)
	if existingRewardCount > 0 {
		l.Logger.Infof("订单 %d 已分配过奖励，跳过", orderID)
		return nil
	}

	// 1. 查询订单信息
	var order model.Order
	if err := l.svcCtx.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		return fmt.Errorf("查询订单信息失败: %w", err)
	}

	// 2. 查询活动信息，检查是否启用分销
	var campaign model.Campaign
	if err := l.svcCtx.DB.Where("id = ?", campaignID).First(&campaign).Error; err != nil {
		return fmt.Errorf("查询活动信息失败: %w", err)
	}

	// 检查活动是否启用分销
	if !campaign.EnableDistribution {
		// 活动未启用分销，使用原有单级奖励逻辑
		return l.calculateSingleLevelReward(&order, referrerID)
	}

	// 3. 解析分销奖励配置
	var rewardConfig DistributionRewardConfig
	if campaign.DistributionRewards != nil {
		if err := json.Unmarshal([]byte(*campaign.DistributionRewards), &rewardConfig); err != nil {
			return fmt.Errorf("解析分销奖励配置失败: %w", err)
		}
	} else {
		// 使用默认配置
		rewardConfig = DistributionRewardConfig{
			Level1: 10,
			Level2: 0,
			Level3: 0,
		}
	}

	// 4. 查询分销链
	distributorPath := order.DistributorPath
	if distributorPath == "" && referrerID > 0 {
		// 如果订单没有分销链，但有推荐人，查询推荐人是否是分销商
		var referrerDistributor model.Distributor
		err := l.svcCtx.DB.Where("user_id = ? AND brand_id = ? AND status = ?", referrerID, order.CampaignId, "active").
			First(&referrerDistributor).Error
		if err == nil {
			// 推荐人是分销商，计算分销链
			autoUpgradeLogic := NewAutoUpgradeLogic(l.ctx, l.svcCtx)
			distributorPath = autoUpgradeLogic.CalculateDistributorPath(referrerDistributor.Id, int(campaign.DistributionLevel))
		}
	}

	// 5. 如果没有分销链，使用原有单级奖励逻辑
	if distributorPath == "" {
		return l.calculateSingleLevelReward(&order, referrerID)
	}

	// 6. 开启事务处理奖励分配
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 事务内幂等性检查（双重检查）
	var count int64
	tx.Model(&model.DistributorReward{}).Where("order_id = ?", orderID).Count(&count)
	if count > 0 {
		tx.Rollback()
		l.Logger.Infof("订单 %d 已分配过奖励（事务内检查），跳过", orderID)
		return nil
	}

	// 7. 解析分销链并分配奖励
	distributorIDs := l.parseDistributorPath(distributorPath)
	maxLevel := int(campaign.DistributionLevel)

	rewardsCreated := 0
	for i, distributorID := range distributorIDs {
		if i >= maxLevel {
			// 超过最大层级，停止分配
			break
		}

		// 查询分销商信息
		var distributor model.Distributor
		err := tx.Where("id = ? AND status = ?", distributorID, "active").First(&distributor).Error
		if err != nil {
			// 分销商不存在或状态不是active，跳过
			continue
		}

		// 计算奖励比例
		var rewardPercentage float64
		level := i + 1 // i从0开始，level从1开始
		switch level {
		case 1:
			rewardPercentage = rewardConfig.Level1
		case 2:
			rewardPercentage = rewardConfig.Level2
		case 3:
			rewardPercentage = rewardConfig.Level3
		default:
			continue
		}

		if rewardPercentage <= 0 {
			continue // 该级别没有奖励
		}

		// 计算奖励金额
		rewardAmount := orderAmount * rewardPercentage / 100

		// 分配奖励
		if err := l.distributeRewardWithTx(tx, &distributor, &order, level, rewardPercentage, rewardAmount); err != nil {
			l.Errorf("分配奖励失败: distributorID=%d, level=%d, error=%v", distributorID, level, err)
			// 单个奖励失败不影响其他奖励分配，继续
		} else {
			rewardsCreated++
		}
	}

	// 如果至少创建了一条奖励记录，提交事务
	if rewardsCreated > 0 {
		if err := tx.Commit().Error; err != nil {
			l.Errorf("提交奖励分配事务失败: %v", err)
			return fmt.Errorf("提交奖励分配事务失败: %w", err)
		}
		l.Logger.Infof("订单 %d 成功分配 %d 条奖励记录", orderID, rewardsCreated)
	} else {
		tx.Rollback()
		l.Logger.Infof("订单 %d 没有分配任何奖励", orderID)
	}

	return nil
}

// calculateSingleLevelReward 计算单级奖励（原有逻辑）
func (l *MultiLevelRewardLogic) calculateSingleLevelReward(order *model.Order, referrerID int64) error {
	if referrerID == 0 {
		// 没有推荐人，无需分配奖励
		return nil
	}

	// 查询活动奖励规则
	var campaign model.Campaign
	if err := l.svcCtx.DB.Where("id = ?", order.CampaignId).First(&campaign).Error; err != nil {
		return fmt.Errorf("查询活动信息失败: %w", err)
	}

	// 使用reward_rule作为奖励金额
	rewardAmount := campaign.RewardRule
	if rewardAmount <= 0 {
		return nil
	}

	// 查询推荐人是否是分销商
	var referrerDistributor model.Distributor
	err := l.svcCtx.DB.Where("user_id = ? AND brand_id = ? AND status = ?", referrerID, order.CampaignId, "active").
		First(&referrerDistributor).Error
	if err != nil {
		// 推荐人不是分销商或状态不是active，不分配奖励
		return nil
	}

	// 幂等性检查：使用事务确保原子性
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 事务内幂等性检查
	var count int64
	tx.Model(&model.DistributorReward{}).Where("order_id = ?", order.Id).Count(&count)
	if count > 0 {
		tx.Rollback()
		l.Logger.Infof("订单 %d 已分配过奖励（单级检查），跳过", order.Id)
		return nil
	}

	// 分配奖励
	if err := l.distributeRewardWithTx(tx, &referrerDistributor, order, 1, 0, rewardAmount); err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("提交单级奖励事务失败: %v", err)
		return fmt.Errorf("提交事务失败: %w", err)
	}

	l.Logger.Infof("订单 %d 单级奖励分配成功，金额: %.2f", order.Id, rewardAmount)
	return nil
}

// distributeReward 分配奖励给单个分销商（独立事务）
func (l *MultiLevelRewardLogic) distributeReward(distributor *model.Distributor, order *model.Order, level int, rewardPercentage float64, rewardAmount float64) error {
	// 开启事务
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 创建奖励记录
	reward := &model.DistributorReward{
		DistributorId: distributor.Id,
		UserId:        distributor.UserId,
		OrderId:       order.Id,
		CampaignId:    order.CampaignId,
		Amount:        rewardAmount,
		Level:         level,
		RewardRate:    rewardPercentage,
		Status:        "settled",
		SettledAt:     order.PaidAt,
	}

	if err := tx.Create(reward).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("创建奖励记录失败: %w", err)
	}

	// 2. 更新用户余额（乐观锁）
	if err := l.updateUserBalance(tx, distributor.UserId, rewardAmount); err != nil {
		tx.Rollback()
		return fmt.Errorf("更新用户余额失败: %w", err)
	}

	// 3. 更新分销商累计收益
	if err := tx.Model(&model.Distributor{}).Where("id = ?", distributor.Id).
		Update("total_earnings", gorm.Expr("total_earnings + ?", rewardAmount)).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新分销商累计收益失败: %w", err)
	}

	// 4. 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

// distributeRewardWithTx 分配奖励给单个分销商（使用传入的事务）
func (l *MultiLevelRewardLogic) distributeRewardWithTx(tx *gorm.DB, distributor *model.Distributor, order *model.Order, level int, rewardPercentage float64, rewardAmount float64) error {
	// 1. 创建奖励记录
	reward := &model.DistributorReward{
		DistributorId: distributor.Id,
		UserId:        distributor.UserId,
		OrderId:       order.Id,
		CampaignId:    order.CampaignId,
		Amount:        rewardAmount,
		Level:         level,
		RewardRate:    rewardPercentage,
		Status:        "settled",
		SettledAt:     order.PaidAt,
	}

	if err := tx.Create(reward).Error; err != nil {
		return fmt.Errorf("创建奖励记录失败: %w", err)
	}

	// 2. 更新用户余额（乐观锁）
	if err := l.updateUserBalance(tx, distributor.UserId, rewardAmount); err != nil {
		return fmt.Errorf("更新用户余额失败: %w", err)
	}

	// 3. 更新分销商累计收益
	if err := tx.Model(&model.Distributor{}).Where("id = ?", distributor.Id).
		Update("total_earnings", gorm.Expr("total_earnings + ?", rewardAmount)).Error; err != nil {
		return fmt.Errorf("更新分销商累计收益失败: %w", err)
	}

	return nil
}

// updateUserBalance 更新用户余额（乐观锁）
func (l *MultiLevelRewardLogic) updateUserBalance(tx *gorm.DB, userID int64, amount float64) error {
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		// 查询当前余额和版本号
		var userBalance model.UserBalance
		err := tx.Where("user_id = ?", userID).First(&userBalance).Error
		if err != nil {
			// 如果余额记录不存在，创建一个
			if errors.Is(err, gorm.ErrRecordNotFound) {
				userBalance = model.UserBalance{
					UserId:      userID,
					Balance:     amount,
					TotalReward: amount,
					Version:     0,
				}
				if err := tx.Create(&userBalance).Error; err != nil {
					return fmt.Errorf("创建用户余额失败: %w", err)
				}
				return nil
			}
			return fmt.Errorf("查询用户余额失败: %w", err)
		}

		// 尝试更新（带版本号检查）
		result := tx.Model(&model.UserBalance{}).
			Where("user_id = ? AND version = ?", userID, userBalance.Version).
			Updates(map[string]interface{}{
				"balance":      userBalance.Balance + amount,
				"total_reward": userBalance.TotalReward + amount,
				"version":      userBalance.Version + 1,
			})

		if result.RowsAffected > 0 {
			return nil // 更新成功
		}

		// 版本冲突，等待后重试
		// 在实际应用中，可以添加延时
	}

	return errors.New("更新用户余额失败：达到最大重试次数")
}

// parseDistributorPath 解析分销链路径
func (l *MultiLevelRewardLogic) parseDistributorPath(path string) []int64 {
	var distributorIDs []int64

	// 尝试解析为JSON数组
	if strings.HasPrefix(path, "[") {
		var ids []int64
		if err := json.Unmarshal([]byte(path), &ids); err == nil {
			return ids
		}
	}

	// 尝试解析为逗号分隔的字符串
	if path != "" {
		parts := strings.Split(path, ",")
		for _, part := range parts {
			id, err := strconv.ParseInt(strings.TrimSpace(part), 10, 64)
			if err == nil && id > 0 {
				distributorIDs = append(distributorIDs, id)
			}
		}
	}

	return distributorIDs
}
