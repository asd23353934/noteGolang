package redis

import (
	"NoteGolang/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisDB struct {
	client *redis.Client
}

type RedisConfig struct {
	Host            string
	Port            int
	Password        string
	DB              int
	PoolSize        int
	MinIdleConns    int
	MaxRetries      int
	MaxRetryBackoff time.Duration
}

func NewRedisClient(cfg *config.RedisConfig) (*RedisDB, error) {
	fmt.Println(cfg.RedisHost)
	fmt.Println(cfg.RedisPort)
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
		PoolSize: cfg.RedisPoolSize,
		// MinIdleConns:    cfg.MinIdleConns,
		// MaxRetries:      cfg.MaxRetries,
		// MaxRetryBackoff: cfg.MaxRetryBackoff,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := client.Ping(ctx)
	fmt.Println("redis ping:", result.Val())
	if result.Val() != "PONG" {
		return nil, fmt.Errorf(result.Val())
	}
	return &RedisDB{client: client}, nil
}

func (r *RedisDB) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisDB) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}
