package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPerformanceMiddlewareHandlePassThrough(t *testing.T) {
	m := NewPerformanceMiddleware()
	called := false

	h := m.Handle(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusCreated)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/demo", nil)
	resp := httptest.NewRecorder()
	h(resp, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestResponseWriterWriteHeader(t *testing.T) {
	resp := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: resp, statusCode: http.StatusOK}
	rw.WriteHeader(http.StatusBadRequest)
	assert.Equal(t, http.StatusBadRequest, rw.statusCode)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
