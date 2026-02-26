package syncadapter

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var (
	redisAvailableOnce sync.Once
	redisIsAvailable   bool
)

// checkRedisAvailable 快速检查 Redis 是否可用（超时 500ms）
// 使用 sync.Once 确保只检查一次
func checkRedisAvailable() bool {
	redisAvailableOnce.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:         "localhost:6379",
			Password:     "",
			DB:           0,
			DialTimeout:  500 * time.Millisecond,
			ReadTimeout:  500 * time.Millisecond,
			WriteTimeout: 500 * time.Millisecond,
		})
		defer client.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		_, err := client.Ping(ctx).Result()
		redisIsAvailable = err == nil
	})
	return redisIsAvailable
}

// skipIfNoRedis 跳过 Redis 不可用的测试
func skipIfNoRedis(t *testing.T) *redis.Client {
	if !checkRedisAvailable() {
		t.Skip("Redis not available")
		return nil
	}
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func TestNewSyncQueue(t *testing.T) {
	redisClient := skipIfNoRedis(t)
	if redisClient == nil {
		return
	}
	defer redisClient.Close()

	queue := NewSyncQueue(redisClient, "test_queue")

	assert.NotNil(t, queue)
	assert.Equal(t, "test_queue", queue.key)
	assert.NotNil(t, queue.redis)
}

func TestSyncQueue_Enqueue(t *testing.T) {
	redisClient := skipIfNoRedis(t)
	if redisClient == nil {
		return
	}
	defer redisClient.Close()

	queue := NewSyncQueue(redisClient, "test_enqueue")
	queue.Clear()

	task := &SyncTask{
		TaskId:    "task_1",
		Type:      "order",
		OrderId:   123,
		Attempts:  0,
		CreatedAt: time.Now(),
	}

	err := queue.Enqueue(task)
	assert.NoError(t, err)

	length, err := queue.Length()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), length)
}

func TestSyncQueue_Dequeue(t *testing.T) {
	redisClient := skipIfNoRedis(t)
	if redisClient == nil {
		return
	}
	defer redisClient.Close()

	queue := NewSyncQueue(redisClient, "test_dequeue")

	task := &SyncTask{
		TaskId:    "task_1",
		Type:      "order",
		OrderId:   123,
		Attempts:  0,
		CreatedAt: time.Now(),
	}

	err := queue.Enqueue(task)
	assert.NoError(t, err)

	dequeuedTask, err := queue.Dequeue(5 * time.Second)
	assert.NoError(t, err)
	assert.NotNil(t, dequeuedTask)
	assert.Equal(t, "task_1", dequeuedTask.TaskId)
	assert.Equal(t, "order", dequeuedTask.Type)
	assert.Equal(t, int64(123), dequeuedTask.OrderId)
}

func TestSyncQueue_Dequeue_Timeout(t *testing.T) {
	redisClient := skipIfNoRedis(t)
	if redisClient == nil {
		return
	}
	defer redisClient.Close()

	queue := NewSyncQueue(redisClient, "test_dequeue_timeout")

	err := queue.Clear()
	assert.NoError(t, err)

	task, err := queue.Dequeue(2 * time.Second)
	if err == redis.Nil {
		return
	}
	assert.NoError(t, err)
	assert.Nil(t, task)
}

func TestSyncQueue_Length(t *testing.T) {
	redisClient := skipIfNoRedis(t)
	if redisClient == nil {
		return
	}
	defer redisClient.Close()

	queue := NewSyncQueue(redisClient, "test_length")

	err := queue.Clear()
	assert.NoError(t, err)

	length, err := queue.Length()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), length)

	queue.Enqueue(&SyncTask{TaskId: "1", Type: "order"})
	queue.Enqueue(&SyncTask{TaskId: "2", Type: "order"})

	length, _ = queue.Length()
	assert.Equal(t, int64(2), length)
}

func TestSyncQueue_Clear(t *testing.T) {
	redisClient := skipIfNoRedis(t)
	if redisClient == nil {
		return
	}
	defer redisClient.Close()

	queue := NewSyncQueue(redisClient, "test_clear")

	queue.Enqueue(&SyncTask{TaskId: "1", Type: "order"})
	queue.Enqueue(&SyncTask{TaskId: "2", Type: "order"})

	err := queue.Clear()
	assert.NoError(t, err)

	length, _ := queue.Length()
	assert.Equal(t, int64(0), length)
}
