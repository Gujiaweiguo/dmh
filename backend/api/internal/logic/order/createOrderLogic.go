// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderReq) (resp *types.OrderResp, err error) {
	l.Infof("CreateOrder called: CampaignID=%d, Phone=%s", req.CampaignId, req.Phone)

	if err := validatePhone(req.Phone); err != nil {
		l.Errorf("Invalid phone format: %v", err)
		return nil, fmt.Errorf("手机号格式错误: %v", err)
	}

	if _, err := l.validateCampaign(req.CampaignId); err != nil {
		l.Errorf("Campaign validation failed: %v", err)
		return nil, err
	}

	if err := l.checkDuplicate(req.CampaignId, req.Phone); err != nil {
		l.Errorf("Duplicate order check failed: %v", err)
		return nil, err
	}

	formFields, err := l.getFormFields(req.CampaignId)
	if err != nil {
		l.Errorf("Failed to get form fields: %v", err)
		return nil, fmt.Errorf("获取表单字段配置失败: %v", err)
	}

	if err := l.validateFormData(req.FormData, formFields); err != nil {
		l.Errorf("Form data validation failed: %v", err)
		return nil, fmt.Errorf("表单数据验证失败: %v", err)
	}

	formDataJSON, err := json.Marshal(req.FormData)
	if err != nil {
		l.Errorf("Failed to marshal form data: %v", err)
		return nil, fmt.Errorf("表单数据序列化失败: %v", err)
	}

	timestamp := time.Now().Unix()

	order := &model.Order{
		CampaignId:         req.CampaignId,
		Phone:              req.Phone,
		FormData:           string(formDataJSON),
		ReferrerId:         req.ReferrerId,
		Status:             "pending",
		PayStatus:          "unpaid",
		VerificationStatus: "unverified",
	}

	if err := l.svcCtx.DB.Create(order).Error; err != nil {
		l.Errorf("Failed to create order: %v", err)
		if isDuplicateOrderError(err) {
			return nil, fmt.Errorf("该手机号已参与此活动，请勿重复报名")
		}
		return nil, fmt.Errorf("创建订单失败: %v", err)
	}

	verificationCode := l.generateVerificationCode(order.Id, req.Phone, timestamp)
	if err := l.svcCtx.DB.Model(order).Update("verification_code", verificationCode).Error; err != nil {
		l.Errorf("Failed to update verification code: %v", err)
		return nil, fmt.Errorf("更新核销码失败: %v", err)
	}
	order.VerificationCode = verificationCode

	l.Infof("Order created successfully: ID=%d, CampaignID=%d, Phone=%s", order.Id, order.CampaignId, order.Phone)

	resp = &types.OrderResp{
		Id:         order.Id,
		CampaignId: order.CampaignId,
		Phone:      order.Phone,
		FormData:   req.FormData,
		ReferrerId: order.ReferrerId,
		Status:     order.Status,
		Amount:     order.Amount,
		CreatedAt:  order.CreatedAt.Format("2006-01-02T15:04:05"),
	}

	l.Infof("Returning response: %+v", resp)
	return resp, nil
}

func (l *CreateOrderLogic) validateCampaign(campaignId int64) (*model.Campaign, error) {
	var campaign model.Campaign
	if err := l.svcCtx.DB.Where("id = ? AND deleted_at IS NULL", campaignId).First(&campaign).Error; err != nil {
		if errors.Is(err, l.svcCtx.DB.Error) {
			return nil, fmt.Errorf("活动不存在")
		}
		l.Errorf("Failed to query campaign: %v", err)
		return nil, fmt.Errorf("查询活动失败: %v", err)
	}

	if campaign.Status != "active" {
		return nil, fmt.Errorf("活动未开始或已结束")
	}

	now := time.Now()
	if now.Before(campaign.StartTime) {
		return nil, fmt.Errorf("活动尚未开始")
	}
	if now.After(campaign.EndTime) {
		return nil, fmt.Errorf("活动已结束")
	}

	return &campaign, nil
}

func (l *CreateOrderLogic) checkDuplicate(campaignId int64, phone string) error {
	var count int64
	if err := l.svcCtx.DB.Model(&model.Order{}).
		Where("campaign_id = ? AND phone = ? AND deleted_at IS NULL", campaignId, phone).
		Count(&count).Error; err != nil {
		l.Errorf("Failed to check duplicate order: %v", err)
		return fmt.Errorf("检查重复订单失败: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("该手机号已参与此活动，请勿重复报名")
	}

	return nil
}

func (l *CreateOrderLogic) getFormFields(campaignId int64) ([]model.FormField, error) {
	var campaign model.Campaign
	if err := l.svcCtx.DB.Where("id = ? AND deleted_at IS NULL", campaignId).First(&campaign).Error; err != nil {
		return nil, fmt.Errorf("查询活动失败: %v", err)
	}

	var formFields []model.FormField
	if err := json.Unmarshal([]byte(campaign.FormFields), &formFields); err != nil {
		l.Errorf("Failed to parse form fields: %v", err)
		return nil, fmt.Errorf("解析表单字段配置失败: %v", err)
	}

	return formFields, nil
}

func (l *CreateOrderLogic) validateFormData(formData map[string]string, formFields []model.FormField) error {
	for _, field := range formFields {
		if field.Required {
			value, exists := formData[field.Name]
			if !exists || strings.TrimSpace(value) == "" {
				return fmt.Errorf("必填字段 %s 不能为空", field.Label)
			}

			if err := validateField(value, field); err != nil {
				return fmt.Errorf("字段 %s 验证失败: %v", field.Label, err)
			}
		}
	}

	return nil
}

func validateField(value string, field model.FormField) error {
	switch field.Type {
	case "text":
		return validateText(value, field)
	case "phone":
		return validatePhone(value)
	case "email":
		return validateEmail(value)
	case "number":
		return validateNumber(value)
	case "textarea":
		return validateTextarea(value)
	case "address":
		return validateAddress(value)
	case "select":
		return validateSelect(value, field)
	default:
		return fmt.Errorf("不支持的字段类型: %s", field.Type)
	}
}

func validateText(value string, field model.FormField) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("文本不能为空")
	}
	return nil
}

func validatePhone(value string) error {
	pattern := `^1[3-9]\d{9}$`
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("手机号格式不正确")
	}
	return nil
}

func validateEmail(value string) error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("邮箱格式不正确")
	}
	return nil
}

func validateNumber(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("数字不能为空")
	}
	return nil
}

func validateTextarea(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("内容不能为空")
	}
	return nil
}

func validateAddress(value string) error {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) < 10 {
		return errors.New("地址长度不能少于10个字符")
	}
	if len(trimmed) > 200 {
		return errors.New("地址长度不能超过200个字符")
	}
	return nil
}

func validateSelect(value string, field model.FormField) error {
	if len(field.Options) == 0 {
		return errors.New("select类型字段必须配置选项")
	}

	for _, option := range field.Options {
		if value == option {
			return nil
		}
	}

	return fmt.Errorf("请选择有效的选项")
}

func (l *CreateOrderLogic) generateVerificationCode(orderId int64, phone string, timestamp int64) string {
	secretKey := "dmh-verification-secret-2026"

	signatureData := fmt.Sprintf("%d_%s_%d_%s", orderId, phone, timestamp, secretKey)
	hash := md5.Sum([]byte(signatureData))
	signature := hex.EncodeToString(hash[:])

	return fmt.Sprintf("%d_%s_%d_%s", orderId, phone, timestamp, signature)
}

func isDuplicateOrderError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	var mysqlErr *mysqlDriver.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return true
	}

	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "duplicate entry") || strings.Contains(errMsg, "unique constraint failed")
}
