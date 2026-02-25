package ai

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCopywritingHandler(t *testing.T) {
	svcCtx := &svc.ServiceContext{}
	handler := GenerateCopywritingHandler(svcCtx)

	body := types.GenerateCopywritingReq{
		Topic:  "双十一大促",
		Style:  "professional",
		Length: "medium",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/ai/generate-copywriting", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp types.GenerateCopywritingResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.NotEmpty(t, resp.Content)
	assert.Contains(t, resp.Content, "双十一大促")
}

func TestGenerateCopywritingHandlerWithDifferentStyles(t *testing.T) {
	testCases := []struct {
		name         string
		style        string
		expectPrefix string
	}{
		{"casual style", "casual", "🎉"},
		{"urgent style", "urgent", "⚡"},
		{"emotional style", "emotional", "❤️"},
		{"professional style (default)", "professional", "📢"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svcCtx := &svc.ServiceContext{}
			handler := GenerateCopywritingHandler(svcCtx)

			body := types.GenerateCopywritingReq{
				Topic:  "测试活动",
				Style:  tc.style,
				Length: "short",
			}
			bodyBytes, _ := json.Marshal(body)

			req := httptest.NewRequest(http.MethodPost, "/ai/generate-copywriting", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)

			var resp types.GenerateCopywritingResp
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			assert.NotEmpty(t, resp.Content)
			assert.Contains(t, resp.Content, tc.expectPrefix)
		})
	}
}

func TestGenerateCopywritingHandlerWithDifferentLengths(t *testing.T) {
	testCases := []struct {
		name   string
		length string
	}{
		{"short length", "short"},
		{"medium length", "medium"},
		{"long length", "long"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svcCtx := &svc.ServiceContext{}
			handler := GenerateCopywritingHandler(svcCtx)

			body := types.GenerateCopywritingReq{
				Topic:  "测试活动",
				Style:  "professional",
				Length: tc.length,
			}
			bodyBytes, _ := json.Marshal(body)

			req := httptest.NewRequest(http.MethodPost, "/ai/generate-copywriting", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)

			var resp types.GenerateCopywritingResp
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			assert.NotEmpty(t, resp.Content)
		})
	}
}

func TestGenerateCopywritingHandlerDefaultValues(t *testing.T) {
	svcCtx := &svc.ServiceContext{}
	handler := GenerateCopywritingHandler(svcCtx)

	body := types.GenerateCopywritingReq{
		Topic: "默认风格测试",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/ai/generate-copywriting", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp types.GenerateCopywritingResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.NotEmpty(t, resp.Content)
	assert.Contains(t, resp.Content, "📢")
}

func TestGenerateCopywritingHandlerMissingTopic(t *testing.T) {
	svcCtx := &svc.ServiceContext{}
	handler := GenerateCopywritingHandler(svcCtx)

	body := types.GenerateCopywritingReq{}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/ai/generate-copywriting", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGenerateCopywritingHandler_InvalidJSON(t *testing.T) {
	svcCtx := &svc.ServiceContext{}
	handler := GenerateCopywritingHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/ai/generate-copywriting", bytes.NewReader([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler(rec, req)

	// Should return error for invalid JSON
	assert.NotEqual(t, http.StatusOK, rec.Code)
}
