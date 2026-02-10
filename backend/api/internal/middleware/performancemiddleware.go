package middleware

import (
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// PerformanceMiddleware 性能监控中间件
type PerformanceMiddleware struct {
	logger logx.Logger
}

// NewPerformanceMiddleware 创建性能监控中间件
func NewPerformanceMiddleware() *PerformanceMiddleware {
	return &PerformanceMiddleware{
		logger: logx.WithContext(nil),
	}
}

// Handle 处理请求并记录性能指标
func (m *PerformanceMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path
		method := r.Method

		// 包装 ResponseWriter 以获取状态码
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next(wrapped, r)

		duration := time.Since(start)

		// 记录慢请求 (超过 500ms)
		if duration > 500*time.Millisecond {
			m.logger.Slowf("[SLOW] %s %s - %d - %v", method, path, wrapped.statusCode, duration)
		}

		// 记录错误请求 (状态码 >= 400)
		if wrapped.statusCode >= 400 {
			m.logger.Errorf("[ERROR] %s %s - %d - %v", method, path, wrapped.statusCode, duration)
		}

		// 记录所有请求 (DEBUG 级别)
		m.logger.Infof("[REQUEST] %s %s - %d - %v", method, path, wrapped.statusCode, duration)
	}
}

// responseWriter 包装 http.ResponseWriter 以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
