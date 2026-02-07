// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ScanOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewScanOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScanOrderLogic {
	return &ScanOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ScanOrderLogic) ScanOrder(req *types.ScanOrderReq) (resp *types.ScanOrderResp, err error) {
	orderId, phone, timestamp, signature := l.parseVerificationCode(req.Code)
	if orderId == 0 {
		l.Errorf("核销码格式无效: %s", req.Code)
		return nil, fmt.Errorf("核销码无效")
	}

	if !l.verifySignature(req.Code, orderId, phone, timestamp, signature) {
		l.Errorf("核销码签名验证失败: orderId=%d, code=%s", orderId, req.Code)
		return nil, fmt.Errorf("核销码无效")
	}

	var order model.Order
	if err := l.svcCtx.DB.Where("id = ? AND deleted_at IS NULL", orderId).First(&order).Error; err != nil {
		l.Errorf("查询订单失败: %v", err)
		return nil, fmt.Errorf("订单不存在")
	}

	var formData map[string]string
	if err := json.Unmarshal([]byte(order.FormData), &formData); err != nil {
		formData = make(map[string]string)
	}

	l.Infof("订单扫码成功: orderId=%d, code=%s", orderId, req.Code)

	return &types.ScanOrderResp{
		OrderId:   order.Id,
		Status:    order.VerificationStatus,
		PayStatus: order.PayStatus,
		MemberId:  order.MemberID,
		Phone:     order.Phone,
		FormData:  formData,
	}, nil
}

func (l *ScanOrderLogic) parseVerificationCode(code string) (orderId int64, phone string, timestamp int64, signature string) {
	if code == "" {
		return 0, "", 0, ""
	}

	parts := strings.Split(code, "_")
	if len(parts) != 4 {
		return 0, "", 0, ""
	}

	orderIdVal := int64(0)
	if parts[0] != "" {
		orderIdVal, _ = strconv.ParseInt(parts[0], 10, 64)
	}

	timestampVal := int64(0)
	if parts[2] != "" {
		timestampVal, _ = strconv.ParseInt(parts[2], 10, 64)
	}

	return orderIdVal, parts[1], timestampVal, parts[3]
}

func (l *ScanOrderLogic) verifySignature(code string, orderId int64, phone string, timestamp int64, signature string) bool {
	secretKey := "dmh-verification-secret-2026"

	signatureData := fmt.Sprintf("%d_%s_%d_%s", orderId, phone, timestamp, secretKey)
	hash := md5.Sum([]byte(signatureData))
	expectedSignature := hex.EncodeToString(hash[:])

	return signature == expectedSignature
}
