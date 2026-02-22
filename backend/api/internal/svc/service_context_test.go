package svc

import (
	"context"
	"testing"
	"time"

	"dmh/api/internal/config"

	"github.com/redis/go-redis/v9"
)

func TestServiceContextZeroValue(t *testing.T) {
	var s ServiceContext
	if s.DB != nil {
		t.Fatalf("expected nil db on zero value")
	}
}

func TestRedisAdapterMethods(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	t.Cleanup(func() { _ = client.Close() })

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("redis not available: %v", err)
	}

	adapter := &redisAdapter{client: client}
	key := "svc_test_key"

	count, err := adapter.Incr(ctx, key)
	if err != nil {
		t.Fatalf("incr failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("unexpected incr count: %d", count)
	}

	value, err := adapter.Get(ctx, key)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if value != "1" {
		t.Fatalf("unexpected redis value: %s", value)
	}

	if err := adapter.Expire(ctx, key, 60); err != nil {
		t.Fatalf("expire failed: %v", err)
	}

	ttl, err := adapter.TTL(ctx, key)
	if err != nil {
		t.Fatalf("ttl failed: %v", err)
	}
	if ttl == "" {
		t.Fatal("ttl should not be empty")
	}

	_ = client.Del(ctx, key)
}

func TestNewServiceContext(t *testing.T) {
	c := config.Config{}
	c.Mysql.DataSource = "root:Admin168@tcp(127.0.0.1:3306)/dmh_test?charset=utf8mb4&parseTime=true&loc=Local"
	c.Redis.Host = "127.0.0.1:6379"
	c.RateLimit.PosterGenerate.Storage = "redis"
	c.RateLimit.PosterGenerate.MaxRequests = 5
	c.RateLimit.PosterGenerate.WindowDuration = 60
	c.RateLimit.Default.Storage = "redis"
	c.RateLimit.Default.MaxRequests = 100
	c.RateLimit.Default.WindowDuration = 60
	c.WeChatPay.CacheTTL = 60
	c.WeChatPay.HTTPTimeoutMs = 1000
	c.WeChatPay.MockEnabled = true

	s := NewServiceContext(c)
	if s == nil {
		t.Fatal("service context should not be nil")
	}
	if s.PasswordService == nil || s.AuditService == nil || s.SessionService == nil {
		t.Fatal("core services should be initialized")
	}
	if s.PosterRateLimiter == nil || s.DefaultRateLimiter == nil {
		t.Fatal("rate limiters should be initialized")
	}
	if s.PermissionMiddleware == nil {
		t.Fatal("permission middleware should be initialized")
	}

	if s.DB != nil {
		sqlDB, err := s.DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}
}

func TestNewServiceContextWithInvalidDSN(t *testing.T) {
	c := config.Config{}
	c.Mysql.DataSource = "invalid-dsn"
	c.RateLimit.PosterGenerate.Storage = "memory"
	c.RateLimit.PosterGenerate.MaxRequests = 1
	c.RateLimit.PosterGenerate.WindowDuration = 1
	c.RateLimit.Default.Storage = "memory"
	c.RateLimit.Default.MaxRequests = 1
	c.RateLimit.Default.WindowDuration = 1
	c.WeChatPay.CacheTTL = 1
	c.WeChatPay.HTTPTimeoutMs = int((1 * time.Second).Milliseconds())

	s := NewServiceContext(c)
	if s == nil {
		t.Fatal("service context should not be nil even with invalid dsn")
	}
	if s.PasswordService == nil || s.AuditService == nil || s.SessionService == nil {
		t.Fatal("services should still be initialized when db init fails")
	}
}
