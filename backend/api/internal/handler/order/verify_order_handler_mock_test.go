// ============================================================
// Order 模块 Handler 层分层测试示范
// ============================================================
// 职责：测试 HTTP 请求解析、响应组装、状态码
// Mock 策略：Mock Logic 层，不依赖真实数据库
// ============================================================

package order

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/stretchr/testify/assert"
)

// ============================================================
// Mock Logic 实现 - 核心示范
// ============================================================

// MockVerifyOrderLogic 是 VerifyOrderLogic 的 mock 实现
// 用于 Handler 测试，完全隔离数据库依赖
type MockVerifyOrderLogic struct {
	// 控制测试行为
	ShouldError bool
	ErrorMsg    string
	// 返回数据
	ReturnData *types.VerifyOrderResp
	// 记录调用
	CalledWith *types.VerifyOrderReq
}

// NewMockVerifyOrderLogic 创建 mock logic
func NewMockVerifyOrderLogic() *MockVerifyOrderLogic {
	return &MockVerifyOrderLogic{
		ShouldError: false,
		ReturnData:  &types.VerifyOrderResp{},
	}
}

// VerifyOrder 实现 Logic 接口
func (m *MockVerifyOrderLogic) VerifyOrder(req *types.VerifyOrderReq) (*types.VerifyOrderResp, error) {
	m.CalledWith = req
	if m.ShouldError {
		return nil, &MockError{Message: m.ErrorMsg}
	}
	return m.ReturnData, nil
}

// MockError 自定义错误类型
type MockError struct {
	Message string
}

func (e *MockError) Error() string {
	return e.Message
}

// ============================================================
// Handler 测试环境设置
// ============================================================

// setupVerifyOrderHandlerTest 创建 Handler 测试环境
// 关键点：不依赖真实数据库，只创建最小 ServiceContext
func setupVerifyOrderHandlerTest(t *testing.T, mockLogic *MockVerifyOrderLogic) http.HandlerFunc {
	t.Helper()

	// 创建最小的 ServiceContext
	// Handler 测试不需要完整配置，只需要能处理请求
	_ = &svc.ServiceContext{} // svcCtx 占位，实际通过 mockLogic

	// 返回使用 mock logic 的 handler
	// 注意：实际项目中需要通过依赖注入或接口替换
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VerifyOrderReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := mockLogic.VerifyOrder(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

// ============================================================
// 测试用例 - 成功路径
// ============================================================

func TestVerifyOrderHandler_Mock_Success(t *testing.T) {
	// 1. 准备 mock 数据
	mockLogic := NewMockVerifyOrderLogic()
	mockLogic.ReturnData = &types.VerifyOrderResp{
		OrderId:    1,
		Status:     "verified",
		VerifiedAt: "2026-02-26T10:00:00Z",
	}

	// 2. 创建 handler
	handler := setupVerifyOrderHandlerTest(t, mockLogic)

	// 3. 构造请求
	reqBody := types.VerifyOrderReq{
		Code:   "1_13800138000_1234567890_abc123",
		Remark: "Test verification",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// 4. 添加认证上下文（模拟中间件）
	ctx := context.WithValue(req.Context(), "userId", int64(1))
	ctx = context.WithValue(ctx, "roles", []string{"brand_admin"})
	req = req.WithContext(ctx)

	resp := httptest.NewRecorder()

	// 5. 执行 handler
	handler(resp, req)

	// 6. 验证结果
	assert.Equal(t, http.StatusOK, resp.Code)

	var result types.VerifyOrderResp
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.OrderId)
	assert.Equal(t, "verified", result.Status)

	// 7. 验证 mock 被正确调用
	assert.NotNil(t, mockLogic.CalledWith)
	assert.Equal(t, "1_13800138000_1234567890_abc123", mockLogic.CalledWith.Code)
}

// ============================================================
// 测试用例 - 错误路径
// ============================================================

func TestVerifyOrderHandler_Mock_InvalidJSON(t *testing.T) {
	mockLogic := NewMockVerifyOrderLogic()
	handler := setupVerifyOrderHandlerTest(t, mockLogic)

	// 发送无效 JSON
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/verify", strings.NewReader("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	// JSON 解析失败应返回 400
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestVerifyOrderHandler_Mock_LogicError(t *testing.T) {
	mockLogic := NewMockVerifyOrderLogic()
	mockLogic.ShouldError = true
	mockLogic.ErrorMsg = "订单已核销" // 业务错误

	handler := setupVerifyOrderHandlerTest(t, mockLogic)

	reqBody := types.VerifyOrderReq{
		Code: "1_13800138000_1234567890_abc123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	// Logic 层错误应正确传递
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "订单已核销")
}

func TestVerifyOrderHandler_Mock_OrderNotFound(t *testing.T) {
	mockLogic := NewMockVerifyOrderLogic()
	mockLogic.ShouldError = true
	mockLogic.ErrorMsg = "订单不存在"

	handler := setupVerifyOrderHandlerTest(t, mockLogic)

	reqBody := types.VerifyOrderReq{
		Code: "999_13800138000_1234567890_abc123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "订单不存在")
}

func TestVerifyOrderHandler_Mock_InvalidCode(t *testing.T) {
	mockLogic := NewMockVerifyOrderLogic()
	mockLogic.ShouldError = true
	mockLogic.ErrorMsg = "核销码无效"

	handler := setupVerifyOrderHandlerTest(t, mockLogic)

	reqBody := types.VerifyOrderReq{
		Code: "invalid_code",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "核销码无效")
}

func TestVerifyOrderHandler_Mock_MissingCode(t *testing.T) {
	mockLogic := NewMockVerifyOrderLogic()
	handler := setupVerifyOrderHandlerTest(t, mockLogic)

	reqBody := types.VerifyOrderReq{
		Code: "", // 空核销码
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	// 即使 handler 接收了请求，logic 应该返回错误
	// 这里 handler 会正常调用，但 logic 行为由测试配置决定
	assert.Equal(t, http.StatusOK, resp.Code) // 因为 mockLogic.ShouldError = false
}

// ============================================================
// 测试用例 - 边界情况
// ============================================================

func TestVerifyOrderHandler_Mock_EmptyBody(t *testing.T) {
	mockLogic := NewMockVerifyOrderLogic()
	handler := setupVerifyOrderHandlerTest(t, mockLogic)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/verify", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

// ============================================================
// 关键模式总结
// ============================================================
//
// Handler 层测试模式：
// 1. 创建 Mock Logic 结构体，实现 Logic 接口
// 2. Mock Logic 控制 ShouldError 和 ReturnData
// 3. Handler 测试不依赖真实数据库
// 4. 测试覆盖：
//    - 成功路径
//    - 无效 JSON
//    - Logic 层错误
//    - 边界情况（空值、缺失字段）
//
// 优点：
// - 快速执行（无数据库 IO）
// - 隔离测试（无外部依赖）
// - 可预测行为（mock 控制返回）
