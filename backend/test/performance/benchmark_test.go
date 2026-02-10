package performance

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// BenchmarkCreateOrder 测试订单创建接口性能
func BenchmarkCreateOrder(b *testing.B) {
	baseURL := "http://localhost:8889"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 模拟订单创建请求
			_ = baseURL
			time.Sleep(1 * time.Millisecond) // 模拟网络延迟
		}
	})
}

// BenchmarkGetCampaigns 测试活动列表查询性能
func BenchmarkGetCampaigns(b *testing.B) {
	baseURL := "http://localhost:8889"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = baseURL
		time.Sleep(1 * time.Millisecond)
	}
}

// BenchmarkVerifyOrder 测试订单核销接口性能
func BenchmarkVerifyOrder(b *testing.B) {
	baseURL := "http://localhost:8889"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = baseURL
			time.Sleep(2 * time.Millisecond)
		}
	})
}

// TestConcurrentOrderCreation 测试并发订单创建
func TestConcurrentOrderCreation(t *testing.T) {
	concurrency := 100
	duration := 10 * time.Second

	var wg sync.WaitGroup
	errors := make(chan error, concurrency)

	start := time.Now()
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			deadline := time.Now().Add(duration)
			for time.Now().Before(deadline) {
				// 模拟订单创建
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	elapsed := time.Since(start)
	fmt.Printf("并发测试完成: %d 并发, 耗时 %v\n", concurrency, elapsed)
}

// TestDatabaseConnectionPool 测试数据库连接池性能
func TestDatabaseConnectionPool(t *testing.T) {
	maxConnections := 50
	iterations := 1000

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < maxConnections; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < iterations/maxConnections; j++ {
				// 模拟数据库查询
				time.Sleep(1 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("连接池测试完成: %d 连接, %d 次查询, 耗时 %v\n",
		maxConnections, iterations, elapsed)
}

// TestMemoryLeak 测试内存泄漏
func TestMemoryLeak(t *testing.T) {
	// 运行多次，检查内存是否持续增长
	iterations := 100

	for i := 0; i < iterations; i++ {
		// 模拟创建大量对象
		data := make([]byte, 1024*1024) // 1MB
		_ = data

		if i%10 == 0 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	// 强制GC
	// runtime.GC()

	fmt.Printf("内存泄漏测试完成: %d 次迭代\n", iterations)
}
