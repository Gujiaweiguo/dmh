// Redis availability probe for integration tests
package testutil

import (
	"context"
	"net"
	"testing"
	"time"
)

const (
	// RedisProbeTimeout is the maximum time to wait for Redis
	RedisProbeTimeout = 5 * time.Second
)

// ProbeRedis checks if Redis is available at the given address
func ProbeRedis(addr string) error {
	// Use simple TCP dial instead of redis client to avoid dependency
	ctx, cancel := context.WithTimeout(context.Background(), RedisProbeTimeout)
	defer cancel()

	// Simple TCP connection check
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

// SkipIfNoRedis skips the test if Redis is not available
// Only use this for tests that absolutely require Redis
func SkipIfNoRedis(t *testing.T) {
	t.Helper()
	addr := getEnv("REDIS_TEST_HOST", "localhost:6379")

	if err := ProbeRedis(addr); err != nil {
		t.Skipf("SKIP_REASON:%s | Redis at %s not available: %v",
			SkipReasonRedisUnavailable, addr, err)
	}
}

// IsRedisAvailable checks if Redis is available without skipping
func IsRedisAvailable() bool {
	addr := getEnv("REDIS_TEST_HOST", "localhost:6379")
	return ProbeRedis(addr) == nil
}
