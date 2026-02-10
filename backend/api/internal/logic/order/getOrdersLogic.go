package order

import (
	"context"
	"encoding/json"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrdersLogic {
	return &GetOrdersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrdersLogic) GetOrders() (resp *types.OrderListResp, err error) {
	var modelOrders []model.Order

	if err := l.svcCtx.DB.Model(&model.Order{}).Order("created_at DESC").Find(&modelOrders).Error; err != nil {
		l.Errorf("查询订单列表失败: %v", err)
		return nil, err
	}

	orders := make([]types.OrderResp, 0, len(modelOrders))
	for _, order := range modelOrders {
		formData := make(map[string]string)
		if order.FormData != "" {
			if err := json.Unmarshal([]byte(order.FormData), &formData); err != nil {
				l.Errorf("Failed to parse form data for order %d: %v", order.Id, err)
			}
		}

		orders = append(orders, types.OrderResp{
			Id:         order.Id,
			CampaignId: order.CampaignId,
			Phone:      order.Phone,
			FormData:   formData,
			ReferrerId: order.ReferrerId,
			Status:     order.Status,
			Amount:     order.Amount,
			CreatedAt:  order.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	resp = &types.OrderListResp{
		Total:  int64(len(orders)),
		Orders: orders,
	}
	return resp, nil
}
