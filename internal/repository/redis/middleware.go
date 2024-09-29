package redis

import (
	"context"
	"fmt"
	"time"
)

type MiddlewareRepository struct {
	db *RedisDB
}

func NewMiddlewareRepository(db *RedisDB) *MiddlewareRepository {
	return &MiddlewareRepository{db: db}
}

// Allow 检查是否允许请求
func (r *MiddlewareRepository) Allow(key string, burst int64) (bool, error) {
	ctx := context.Background()
	now := time.Now().Unix()
	windowKey := fmt.Sprintf("%s:%d", key, now)

	pipe := r.db.client.Pipeline()
	incr := pipe.Incr(ctx, windowKey)
	pipe.Expire(ctx, windowKey, time.Second)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("Redis pipeline failed: %v", err)
	}

	count, err := incr.Result()
	if err != nil {
		return false, fmt.Errorf("Redis Incr failed: %v", err)
	}

	return count <= int64(burst), nil
}
