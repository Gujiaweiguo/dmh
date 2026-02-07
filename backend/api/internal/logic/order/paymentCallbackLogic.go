// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"
	"errors"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type PaymentCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPaymentCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentCallbackLogic {
	return &PaymentCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentCallbackLogic) PaymentCallback(req *types.PaymentCallbackReq) error {
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			l.Errorf("Payment callback panic recovered: %v", r)
		}
	}()

	var order model.Order
	if err := tx.Where("id = ?", req.OrderId).First(&order).Error; err != nil {
		tx.Rollback()
		l.Errorf("Failed to query order: %v", err)
		return errors.New("Order not found")
	}

	if order.PayStatus == "paid" {
		tx.Rollback()
		l.Infof("Order already paid: orderId=%d", req.OrderId)
		return nil
	}

	order.PayStatus = "paid"
	order.TradeNo = req.TradeNo
	order.Amount = req.Amount
	now := time.Now()
	order.UpdatedAt = now

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		l.Errorf("Failed to update order: %v", err)
		return err
	}

	if err := l.calculateAndSettleRewards(tx, order); err != nil {
		tx.Rollback()
		l.Errorf("Failed to calculate rewards: %v", err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		l.Errorf("Failed to commit transaction: %v", err)
		return err
	}

	l.Infof("Payment callback processed successfully: orderId=%d, tradeNo=%s", req.OrderId, req.TradeNo)

	return nil
}

func (l *PaymentCallbackLogic) calculateAndSettleRewards(tx *gorm.DB, order model.Order) error {
	var campaign model.Campaign
	if err := tx.Where("id = ?", order.CampaignId).First(&campaign).Error; err != nil {
		l.Errorf("Failed to query campaign: %v", err)
		return errors.New("Campaign not found")
	}

	if order.ReferrerId == 0 || !campaign.EnableDistribution {
		l.Infof("No referrer or distribution disabled for order: orderId=%d", order.Id)
		return nil
	}

	var referrerDistributor model.Distributor
	if err := tx.Where("user_id = ? AND status = ?", order.ReferrerId, "active").
		First(&referrerDistributor).Error; err != nil {
		l.Errorf("Failed to query distributor: %v", err)
		return errors.New("Referrer not found as active distributor")
	}

	rewardAmount := order.Amount * campaign.RewardRule

	var distributorLevelRewards []model.DistributorLevelReward
	if err := tx.Where("brand_id = ?", campaign.BrandId).
		Find(&distributorLevelRewards).Error; err != nil {
		l.Errorf("Failed to query level rewards: %v", err)
		return errors.New("Failed to get reward configuration")
	}

	levelRewardPercent := 0.0
	for _, levelReward := range distributorLevelRewards {
		if levelReward.Level == referrerDistributor.Level {
			levelRewardPercent = levelReward.RewardPercentage
			break
		}
	}

	if levelRewardPercent == 0 {
		l.Errorf("No reward config for distributor level: level=%d", referrerDistributor.Level)
		return errors.New("Reward configuration not found for distributor level")
	}

	actualReward := rewardAmount * (levelRewardPercent / 100)

	now := time.Now()
	rewardRecord := model.DistributorReward{
		DistributorId: referrerDistributor.Id,
		UserId:        referrerDistributor.UserId,
		OrderId:       order.Id,
		CampaignId:    order.CampaignId,
		Amount:        actualReward,
		Level:         referrerDistributor.Level,
		RewardRate:    levelRewardPercent,
		Status:        "settled",
		SettledAt:     &now,
	}

	if err := tx.Create(&rewardRecord).Error; err != nil {
		l.Errorf("Failed to create reward record: %v", err)
		return err
	}

	referrerDistributor.TotalEarnings += actualReward
	if err := tx.Save(&referrerDistributor).Error; err != nil {
		l.Errorf("Failed to update distributor earnings: %v", err)
		return err
	}

	l.Infof("Reward settled successfully: distributorId=%d, orderId=%d, amount=%.2f", referrerDistributor.Id, order.Id, actualReward)

	return nil
}
