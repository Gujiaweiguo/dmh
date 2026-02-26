// ============================================================
// Distributor 模块 Handler 层分层测试示范
// ============================================================
// 职责：测试 HTTP 请求解析、响应组装、状态码
// Mock 策略：Mock Logic 层，不依赖真实数据库
// ============================================================

package distributor

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"dmh/api/internal/types"

	"github.com/stretchr/testify/assert"
)

// ============================================================
// Mock Logic 实现 - 核心示范
// ============================================================

// MockDistributorApplyLogic 是 DistributorApplyLogic 的 mock 实现
type MockDistributorApplyLogic struct {
	ShouldError bool
	ErrorMsg    string
	ReturnData  *types.DistributorApplicationResp
	CalledWith  *types.DistributorApplyReq
}

func NewMockDistributorApplyLogic() *MockDistributorApplyLogic {
	return &MockDistributorApplyLogic{
		ShouldError: false,
		ReturnData:  &types.DistributorApplicationResp{},
	}
}

func (m *MockDistributorApplyLogic) DistributorApply(req *types.DistributorApplyReq) (*types.DistributorApplicationResp, error) {
	m.CalledWith = req
	if m.ShouldError {
		return nil, &MockDistributorError{Message: m.ErrorMsg}
	}
	return m.ReturnData, nil
}

// MockDistributorError 自定义错误类型
type MockDistributorError struct {
	Message string
}

func (e *MockDistributorError) Error() string {
	return e.Message
}

// ============================================================
// Handler 测试环境设置
// ============================================================

// setupDistributorApplyHandlerTest 创建 Handler 测试环境
func setupDistributorApplyHandlerTest(t *testing.T, mockLogic *MockDistributorApplyLogic) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DistributorApplyReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := mockLogic.DistributorApply(&req)
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

func TestDistributorApplyHandler_Mock_Success(t *testing.T) {
	// 1. 准备 mock 数据
	mockLogic := NewMockDistributorApplyLogic()
	now := time.Now()
	mockLogic.ReturnData = &types.DistributorApplicationResp{
		Id:        1,
		UserId:    100,
		BrandId:   1,
		Status:    "pending",
		Reason:    "我想成为分销商",
		CreatedAt: now.Format(time.RFC3339),
	}

	// 2. 创建 handler
	handler := setupDistributorApplyHandlerTest(t, mockLogic)

	// 3. 构造请求
	reqBody := types.DistributorApplyReq{
		BrandId: 1,
		Reason:  "我想成为分销商",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/distributors/apply", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// 4. 添加认证上下文
	ctx := context.WithValue(req.Context(), "userId", int64(100))
	req = req.WithContext(ctx)

	resp := httptest.NewRecorder()

	// 5. 执行 handler
	handler(resp, req)

	// 6. 验证结果
	assert.Equal(t, http.StatusOK, resp.Code)

	var result types.DistributorApplicationResp
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.Id)
	assert.Equal(t, "pending", result.Status)

	// 7. 验证 mock 被正确调用
	assert.NotNil(t, mockLogic.CalledWith)
	assert.Equal(t, int64(1), mockLogic.CalledWith.BrandId)
	assert.Equal(t, "我想成为分销商", mockLogic.CalledWith.Reason)
}

// ============================================================
// 测试用例 - 错误路径
// ============================================================

func TestDistributorApplyHandler_Mock_InvalidJSON(t *testing.T) {
	mockLogic := NewMockDistributorApplyLogic()
	handler := setupDistributorApplyHandlerTest(t, mockLogic)

	// 发送无效 JSON
	req := httptest.NewRequest(http.MethodPost, "/api/v1/distributors/apply", strings.NewReader("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	// JSON 解析失败应返回 400
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestDistributorApplyHandler_Mock_LogicError_AlreadyExists(t *testing.T) {
	mockLogic := NewMockDistributorApplyLogic()
	mockLogic.ShouldError = true
	mockLogic.ErrorMsg = "您已提交申请，请勿重复申请"

	handler := setupDistributorApplyHandlerTest(t, mockLogic)

	reqBody := types.DistributorApplyReq{
		BrandId: 1,
		Reason:  "我想成为分销商",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/distributors/apply", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "您已提交申请")
}

func TestDistributorApplyHandler_Mock_LogicError_AlreadyDistributor(t *testing.T) {
	mockLogic := NewMockDistributorApplyLogic()
	mockLogic.ShouldError = true
	mockLogic.ErrorMsg = "您已经是该品牌的分销商"

	handler := setupDistributorApplyHandlerTest(t, mockLogic)

	reqBody := types.DistributorApplyReq{
		BrandId: 1,
		Reason:  "我想成为分销商",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/distributors/apply", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "已经是该品牌的分销商")
}

func TestDistributorApplyHandler_Mock_LogicError_NotLoggedIn(t *testing.T) {
	mockLogic := NewMockDistributorApplyLogic()
	mockLogic.ShouldError = true
	mockLogic.ErrorMsg = "用户未登录"

	handler := setupDistributorApplyHandlerTest(t, mockLogic)

	reqBody := types.DistributorApplyReq{
		BrandId: 1,
		Reason:  "我想成为分销商",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/distributors/apply", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// 不设置 userId
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "用户未登录")
}

func TestDistributorApplyHandler_Mock_LogicError_InvalidBrandId(t *testing.T) {
	mockLogic := NewMockDistributorApplyLogic()
	mockLogic.ShouldError = true
	mockLogic.ErrorMsg = "品牌ID无效"

	handler := setupDistributorApplyHandlerTest(t, mockLogic)

	reqBody := types.DistributorApplyReq{
		BrandId: 0, // 无效的品牌ID
		Reason:  "我想成为分销商",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/distributors/apply", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "品牌ID无效")
}

// ============================================================
// 测试用例 - 边界情况
// ============================================================

func TestDistributorApplyHandler_Mock_EmptyReason(t *testing.T) {
	mockLogic := NewMockDistributorApplyLogic()
	mockLogic.ReturnData = &types.DistributorApplicationResp{
		Id:        1,
		UserId:    100,
		BrandId:   1,
		Status:    "pending",
		Reason:    "",
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	handler := setupDistributorApplyHandlerTest(t, mockLogic)

	reqBody := types.DistributorApplyReq{
		BrandId: 1,
		Reason:  "", // 空原因
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/distributors/apply", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	// 空原因可能是允许的
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestDistributorApplyHandler_Mock_EmptyBody(t *testing.T) {
	mockLogic := NewMockDistributorApplyLogic()
	handler := setupDistributorApplyHandlerTest(t, mockLogic)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/distributors/apply", strings.NewReader(""))
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
//    - Logic 层业务错误（已存在、已申请、未登录等）
//    - 边界情况（空值、缺失字段）
//
// 分销商特有场景：
// - 用户未登录
// - 品牌ID无效
// - 重复申请
// - 已是分销商
