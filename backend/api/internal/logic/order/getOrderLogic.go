// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"
	"encoding/json"
	"fmt"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderLogic) GetOrder(orderId int64) (resp *types.OrderResp, err error) {
	order := &model.Order{}
	if err := l.svcCtx.DB.First(order, orderId).Error; err != nil {
		l.Errorf("Order not found: %v", err)
		return nil, fmt.Errorf("order not found")
	}

	formData := make(map[string]string)
	if order.FormData != "" {
		if err := json.Unmarshal([]byte(order.FormData), &formData); err != nil {
			l.Errorf("Failed to parse form data: %v", err)
		}
	}

	resp = &types.OrderResp{
		Id:         order.Id,
		CampaignId: order.CampaignId,
		Phone:      order.Phone,
		FormData:   formData,
		ReferrerId: order.ReferrerId,
		Status:     order.Status,
		Amount:     order.Amount,
		CreatedAt:  order.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}
