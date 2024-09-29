package repository

import (
	"NoteGolang/internal/repository/cassandra"
	"NoteGolang/internal/repository/redis"
)

type Repositories struct {
	ArticleRepo     *cassandra.ArticleRepository
	ArticleCache    *redis.ArticleRepository
	MiddlewareCache *redis.MiddlewareRepository
}

func NewRepositories(cassandraDB *cassandra.CassandraDB, redisClient *redis.RedisDB) (*Repositories, error) {

	return &Repositories{
		ArticleRepo:     cassandra.NewArticleRepository(cassandraDB),
		ArticleCache:    redis.NewArticleRepository(redisClient),
		MiddlewareCache: redis.NewMiddlewareRepository(redisClient),
	}, nil
}
