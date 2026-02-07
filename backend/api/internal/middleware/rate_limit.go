package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type RateLimiter interface {
	Allow(userID string) bool
	GetRemaining(userID string) int
	GetResetTime(userID string) time.Time
}

type MemoryRateLimiter struct {
	mu       sync.RWMutex
	requests map[string]*userRequests
	max      int
	duration time.Duration
}

type userRequests struct {
	timestamps []time.Time
}

func NewMemoryRateLimiter(maxRequests int, duration time.Duration) *MemoryRateLimiter {
	limiter := &MemoryRateLimiter{
		requests: make(map[string]*userRequests),
		max:      maxRequests,
		duration: duration,
	}

	go limiter.cleanupExpired()

	return limiter
}

func (l *MemoryRateLimiter) Allow(userID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	ur, exists := l.requests[userID]

	if !exists {
		l.requests[userID] = &userRequests{
			timestamps: []time.Time{now},
		}
		return true
	}

	cutoff := now.Add(-l.duration)
	validTimestamps := make([]time.Time, 0)
	for _, ts := range ur.timestamps {
		if ts.After(cutoff) {
			validTimestamps = append(validTimestamps, ts)
		}
	}

	if len(validTimestamps) >= l.max {
		return false
	}

	validTimestamps = append(validTimestamps, now)
	l.requests[userID].timestamps = validTimestamps

	return true
}

func (l *MemoryRateLimiter) GetRemaining(userID string) int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	ur, exists := l.requests[userID]
	if !exists {
		return l.max
	}

	now := time.Now()
	cutoff := now.Add(-l.duration)
	count := 0

	for _, ts := range ur.timestamps {
		if ts.After(cutoff) {
			count++
		}
	}

	return l.max - count
}

func (l *MemoryRateLimiter) GetResetTime(userID string) time.Time {
	l.mu.RLock()
	defer l.mu.RUnlock()

	ur, exists := l.requests[userID]
	if !exists {
		return time.Now()
	}

	if len(ur.timestamps) == 0 {
		return time.Now()
	}

	return ur.timestamps[0].Add(l.duration)
}

func (l *MemoryRateLimiter) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()

		toDelete := make([]string, 0)
		for userID, ur := range l.requests {
			if len(ur.timestamps) > 0 {
				cutoff := time.Now().Add(-l.duration)
				hasValid := false
				for _, ts := range ur.timestamps {
					if ts.After(cutoff) {
						hasValid = true
						break
					}
				}
				if !hasValid {
					toDelete = append(toDelete, userID)
				}
			}
		}

		for _, id := range toDelete {
			delete(l.requests, id)
		}

		l.mu.Unlock()
	}
}

type RedisRateLimiter struct {
	mu          sync.Mutex
	redis       RedisClient
	prefix      string
	maxRequests int
	duration    time.Duration
}

func NewRedisRateLimiter(redis RedisClient, prefix string, maxRequests int, duration time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{
		redis:       redis,
		prefix:      prefix,
		maxRequests: maxRequests,
		duration:    duration,
	}
}

func (l *RedisRateLimiter) Allow(userID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	key := fmt.Sprintf("%s:requests:%s", l.prefix, userID)

	count, err := l.redis.Incr(context.Background(), key)
	if err != nil {
		logx.Errorf("Redis INCR失败: %v", err)
		return true
	}

	if count == 1 {
		l.redis.Expire(context.Background(), key, int(l.duration.Seconds()))
	}

	if count > int64(l.maxRequests) {
		return false
	}

	return true
}

func (l *RedisRateLimiter) GetRemaining(userID string) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	key := fmt.Sprintf("%s:requests:%s", l.prefix, userID)

	count, err := l.redis.Get(context.Background(), key)
	if err != nil {
		logx.Errorf("Redis GET失败: %v", err)
		return l.maxRequests
	}

	countInt, _ := strconv.Atoi(count)
	remaining := l.maxRequests - countInt
	if remaining < 0 {
		return 0
	}

	return remaining
}

func (l *RedisRateLimiter) GetResetTime(userID string) time.Time {
	l.mu.Lock()
	defer l.mu.Unlock()

	key := fmt.Sprintf("%s:requests:expire:%s", l.prefix, userID)

	ttl, err := l.redis.TTL(context.Background(), key)
	if err != nil {
		logx.Errorf("Redis TTL失败: %v", err)
		return time.Now()
	}

	ttlInt, _ := strconv.Atoi(ttl)
	return time.Now().Add(time.Duration(ttlInt) * time.Second)
}

type RedisClient interface {
	Incr(ctx context.Context, key string) (int64, error)
	Get(ctx context.Context, key string) (string, error)
	Expire(ctx context.Context, key string, seconds int) error
	TTL(ctx context.Context, key string) (string, error)
}

func RateLimitMiddleware(limiter RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := getUserID(r)

			if !limiter.Allow(userID) {
				remaining := limiter.GetRemaining(userID)
				resetTime := limiter.GetResetTime(userID)

				logx.WithContext(r.Context()).Infof("Rate limit exceeded for user %s, remaining: %d, reset: %v",
					userID, remaining, resetTime)

				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
				w.Header().Set("X-RateLimit-Reset", resetTime.Format(time.RFC3339))
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(fmt.Sprintf(`{
	"code": 429,
	"msg": "请求过于频繁，请稍后再试",
	"remaining": %d,
	"resetAt": "`+resetTime.Format(time.RFC3339)+`"
		}`, remaining)))
				return
			}

			remaining := limiter.GetRemaining(userID)
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

			next.ServeHTTP(w, r)
		})
	}
}

func getUserID(r *http.Request) string {
	userID := r.Context().Value("userId")
	if userID != nil {
		return fmt.Sprintf("%v", userID)
	}

	ip := r.RemoteAddr
	if ip == "" {
		ip = "unknown"
	}
	return ip
}
