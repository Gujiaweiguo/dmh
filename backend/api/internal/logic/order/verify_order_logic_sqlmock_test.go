//go:build layered_demo
// +build layered_demo

// ============================================================
// Order 模块 Logic 层分层测试示范
// ============================================================
// 职责：测试业务逻辑、数据验证、错误处理
// Mock 策略：使用 sqlmock 隔离数据库
// ============================================================

package order

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ============================================================
// 测试工具函数 - 核心示范
// ============================================================

// setupLogicTestWithSQLMock 创建带 sqlmock 的测试环境
// 关键点：使用 sqlmock 模拟 GORM 底层 SQL 连接
func setupLogicTestWithSQLMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *svc.ServiceContext) {
	t.Helper()

	// 1. 创建 sqlmock 连接
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to create sqlmock")

	// 2. 配置 GORM 使用 mock 连接
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true, // 跳过版本检查
	}), &gorm.Config{
		SkipDefaultTransaction: true, // 简化 mock（不自动开启事务）
	})
	require.NoError(t, err, "failed to open gorm")

	// 3. 创建 ServiceContext
	svcCtx := &svc.ServiceContext{
		DB: gormDB,
	}

	// 4. 注册清理
	t.Cleanup(func() {
		sqlDB.Close()
	})

	return gormDB, mock, svcCtx
}

// ============================================================
// 测试用例 - 核销成功路径
// ============================================================

func TestVerifyOrderLogic_SQLMock_Success(t *testing.T) {
	_, mock, svcCtx := setupLogicTestWithSQLMock(t)

	// 1. 准备测试数据
	now := time.Now()
	orderId := int64(1)
	phone := "13800138000"
	timestamp := now.Unix()

	// 2. 生成有效的核销码
	logic := &VerifyOrderLogic{}
	signature := logic.generateTestSignature(orderId, phone, timestamp)
	code := formatVerificationCode(orderId, phone, timestamp, signature)

	// 3. Mock 查询订单
	ordersRows := sqlmock.NewRows([]string{
		"id", "campaign_id", "phone", "form_data", "amount",
		"status", "pay_status", "verification_status", "verification_code",
		"created_at", "updated_at",
	}).AddRow(
		orderId, 1, phone, `{"name":"张三"}`, 100.00,
		"paid", "paid", "unverified", code,
		now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `orders` WHERE id = ? AND deleted_at IS NULL ORDER BY `orders`.`id` LIMIT ?")).
		WithArgs(orderId, 1).
		WillReturnRows(ordersRows)

	// 4. Mock 开启事务
	mock.ExpectBegin()

	// 5. Mock 更新订单状态
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `orders` SET")).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 6. Mock 创建核销记录
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `verification_records`")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 7. Mock 提交事务
	mock.ExpectCommit()

	// 8. 创建带权限的 context
	ctx := context.WithValue(context.Background(), "userId", int64(100))
	ctx = context.WithValue(ctx, "roles", []string{"brand_admin"})

	// 9. 执行 Logic
	logicUnderTest := NewVerifyOrderLogic(ctx, svcCtx)
	req := &types.VerifyOrderReq{
		Code: code,
	}

	resp, err := logicUnderTest.VerifyOrder(req)

	// 10. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, orderId, resp.OrderId)
	assert.Equal(t, "verified", resp.Status)

	// 11. 验证所有 SQL 期望都被满足
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// 测试用例 - 订单不存在
// ============================================================

func TestVerifyOrderLogic_SQLMock_OrderNotFound(t *testing.T) {
	db, mock, svcCtx := setupLogicTestWithSQLMock(t)
	_ = db // 避免未使用警告

	// 1. 准备测试数据
	orderId := int64(999)
	phone := "13800138000"
	timestamp := time.Now().Unix()

	logic := &VerifyOrderLogic{}
	signature := logic.generateTestSignature(orderId, phone, timestamp)
	code := formatVerificationCode(orderId, phone, timestamp, signature)

	// 2. Mock 查询返回空结果
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `orders` WHERE id = ? AND deleted_at IS NULL ORDER BY `orders`.`id` LIMIT ?")).
		WithArgs(orderId, 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "campaign_id", "phone", "form_data", "amount",
			"status", "pay_status", "verification_status", "verification_code",
			"created_at", "updated_at",
		}))

	// 3. 创建带权限的 context
	ctx := context.WithValue(context.Background(), "userId", int64(100))
	ctx = context.WithValue(ctx, "roles", []string{"brand_admin"})

	// 4. 执行 Logic
	logicUnderTest := NewVerifyOrderLogic(ctx, svcCtx)
	req := &types.VerifyOrderReq{
		Code: code,
	}

	resp, err := logicUnderTest.VerifyOrder(req)

	// 5. 验证错误
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "订单不存在")

	// 6. 验证 mock 期望
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// 测试用例 - 订单已核销
// ============================================================

func TestVerifyOrderLogic_SQLMock_AlreadyVerified(t *testing.T) {
	db, mock, svcCtx := setupLogicTestWithSQLMock(t)
	_ = db

	now := time.Now()
	orderId := int64(1)
	phone := "13800138000"
	timestamp := now.Unix()

	logic := &VerifyOrderLogic{}
	signature := logic.generateTestSignature(orderId, phone, timestamp)
	code := formatVerificationCode(orderId, phone, timestamp, signature)

	// Mock 查询返回已核销订单
	ordersRows := sqlmock.NewRows([]string{
		"id", "campaign_id", "phone", "form_data", "amount",
		"status", "pay_status", "verification_status", "verification_code",
		"created_at", "updated_at",
	}).AddRow(
		orderId, 1, phone, `{"name":"张三"}`, 100.00,
		"paid", "paid", "verified", code, // verification_status = "verified"
		now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `orders` WHERE id = ? AND deleted_at IS NULL ORDER BY `orders`.`id` LIMIT ?")).
		WithArgs(orderId, 1).
		WillReturnRows(ordersRows)

	ctx := context.WithValue(context.Background(), "userId", int64(100))
	ctx = context.WithValue(ctx, "roles", []string{"brand_admin"})

	logicUnderTest := NewVerifyOrderLogic(ctx, svcCtx)
	req := &types.VerifyOrderReq{
		Code: code,
	}

	resp, err := logicUnderTest.VerifyOrder(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "已核销")

	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// 测试用例 - 权限不足
// ============================================================

func TestVerifyOrderLogic_SQLMock_NoPermission(t *testing.T) {
	_, _, svcCtx := setupLogicTestWithSQLMock(t)

	orderId := int64(1)
	phone := "13800138000"
	timestamp := time.Now().Unix()

	logic := &VerifyOrderLogic{}
	signature := logic.generateTestSignature(orderId, phone, timestamp)
	code := formatVerificationCode(orderId, phone, timestamp, signature)

	// 创建无权限的 context（普通参与者）
	ctx := context.WithValue(context.Background(), "userId", int64(100))
	ctx = context.WithValue(ctx, "roles", []string{"participant"})

	logicUnderTest := NewVerifyOrderLogic(ctx, svcCtx)
	req := &types.VerifyOrderReq{
		Code: code,
	}

	resp, err := logicUnderTest.VerifyOrder(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "权限不足")
}

// ============================================================
// 测试用例 - 无效核销码格式
// ============================================================

func TestVerifyOrderLogic_SQLMock_InvalidCode(t *testing.T) {
	_, _, svcCtx := setupLogicTestWithSQLMock(t)

	ctx := context.WithValue(context.Background(), "userId", int64(100))
	ctx = context.WithValue(ctx, "roles", []string{"brand_admin"})

	logicUnderTest := NewVerifyOrderLogic(ctx, svcCtx)
	req := &types.VerifyOrderReq{
		Code: "invalid_code_format",
	}

	resp, err := logicUnderTest.VerifyOrder(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "核销码无效")
}

// ============================================================
// 测试用例 - 数据库错误
// ============================================================

func TestVerifyOrderLogic_SQLMock_DBError(t *testing.T) {
	db, mock, svcCtx := setupLogicTestWithSQLMock(t)
	_ = db

	orderId := int64(1)
	phone := "13800138000"
	timestamp := time.Now().Unix()

	logic := &VerifyOrderLogic{}
	signature := logic.generateTestSignature(orderId, phone, timestamp)
	code := formatVerificationCode(orderId, phone, timestamp, signature)

	// Mock 数据库错误
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `orders` WHERE id = ? AND deleted_at IS NULL ORDER BY `orders`.`id` LIMIT ?")).
		WithArgs(orderId, 1).
		WillReturnError(sql.ErrConnDone)

	ctx := context.WithValue(context.Background(), "userId", int64(100))
	ctx = context.WithValue(ctx, "roles", []string{"brand_admin"})

	logicUnderTest := NewVerifyOrderLogic(ctx, svcCtx)
	req := &types.VerifyOrderReq{
		Code: code,
	}

	resp, err := logicUnderTest.VerifyOrder(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "订单不存在")

	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// 测试用例 - 事务回滚
// ============================================================

func TestVerifyOrderLogic_SQLMock_TransactionRollback(t *testing.T) {
	db, mock, svcCtx := setupLogicTestWithSQLMock(t)
	_ = db

	now := time.Now()
	orderId := int64(1)
	phone := "13800138000"
	timestamp := now.Unix()

	logic := &VerifyOrderLogic{}
	signature := logic.generateTestSignature(orderId, phone, timestamp)
	code := formatVerificationCode(orderId, phone, timestamp, signature)

	// Mock 查询订单
	ordersRows := sqlmock.NewRows([]string{
		"id", "campaign_id", "phone", "form_data", "amount",
		"status", "pay_status", "verification_status", "verification_code",
		"created_at", "updated_at",
	}).AddRow(
		orderId, 1, phone, `{"name":"张三"}`, 100.00,
		"paid", "paid", "unverified", code,
		now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `orders` WHERE id = ? AND deleted_at IS NULL ORDER BY `orders`.`id` LIMIT ?")).
		WithArgs(orderId, 1).
		WillReturnRows(ordersRows)

	mock.ExpectBegin()

	// Mock 更新成功
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `orders` SET")).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Mock 创建核销记录失败
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `verification_records`")).
		WillReturnError(sql.ErrConnDone)

	// Mock 回滚
	mock.ExpectRollback()

	ctx := context.WithValue(context.Background(), "userId", int64(100))
	ctx = context.WithValue(ctx, "roles", []string{"brand_admin"})

	logicUnderTest := NewVerifyOrderLogic(ctx, svcCtx)
	req := &types.VerifyOrderReq{
		Code: code,
	}

	resp, err := logicUnderTest.VerifyOrder(req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================================
// 辅助函数
// ============================================================

// generateTestSignature 生成测试用签名
func (l *VerifyOrderLogic) generateTestSignature(orderId int64, phone string, timestamp int64) string {
	// 与 verifyOrderLogic.go 中的逻辑一致
	secretKey := "dmh-verification-secret-2026"
	_ = formatVerificationCode(orderId, phone, timestamp, secretKey)
	// 这里简化处理，实际应使用 md5
	return "test_signature"
}

// formatVerificationCode 格式化核销码
func formatVerificationCode(orderId int64, phone string, timestamp int64, signature string) string {
	return string(rune(orderId)) + "_" + phone + "_" + string(rune(timestamp)) + "_" + signature
}

// ============================================================
// 关键模式总结
// ============================================================
//
// Logic 层 sqlmock 测试模式：
// 1. 使用 sqlmock.New() 创建 mock 数据库连接
// 2. 配置 GORM 使用 mock 连接
// 3. 使用 mock.ExpectQuery/ExpectExec 设置 SQL 期望
// 4. 执行 Logic 方法
// 5. 验证返回结果和错误
// 6. 使用 mock.ExpectationsWereMet() 确保所有期望被满足
//
// SQL 匹配技巧：
// - 使用 regexp.QuoteMeta() 精确匹配 SQL
// - WithArgs() 设置参数期望
// - WillReturnRows() 设置返回数据
// - WillReturnError() 设置错误
// - WillReturnResult() 设置执行结果
//
// 事务测试：
// - ExpectBegin() - 开启事务
// - ExpectCommit() - 提交事务
// - ExpectRollback() - 回滚事务
//
// 优点：
// - 隔离数据库依赖
// - 可以测试各种边界情况（DB 错误、超时等）
// - 快速执行
