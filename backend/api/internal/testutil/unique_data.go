package testutil

import (
	"fmt"
	"sync/atomic"
	"time"
)

var uniqueCounter uint64

func nextUnique() uint64 {
	return atomic.AddUint64(&uniqueCounter, 1)
}

func GenUniqueUsername(prefix string) string {
	if prefix == "" {
		prefix = "user"
	}
	return fmt.Sprintf("%s_%d_%d", prefix, time.Now().UnixNano(), nextUnique())
}

func GenUniquePhone() string {
	seed := time.Now().UnixNano() + int64(nextUnique())
	if seed < 0 {
		seed = -seed
	}
	return fmt.Sprintf("138%08d", seed%100000000)
}
