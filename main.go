package main

import (
	"NoteGolang/internal/api/routes"
	"NoteGolang/internal/config"
	"NoteGolang/internal/middleware"
	"NoteGolang/internal/repository"
	"NoteGolang/internal/repository/cassandra"
	"NoteGolang/internal/repository/redis"
	"NoteGolang/internal/service"

	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func init() {

}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cassandraCfg, err := config.CassandraLoad(cfg.Environment)
	if err != nil {
		log.Fatalf("failed cassandraCfg: %w", err)
		return
	}
	cassandraHosts := strings.Split(cassandraCfg.CassandraHost, ",")
	cassandraDB, err := cassandra.NewDB(cassandraHosts, cassandraCfg.CassandraUsername, cassandraCfg.CassandraPassword, cassandraCfg.CassandraPort, "")
	if err != nil {
		log.Fatalf("failed to connect to Cassandra: %w", err)
		return
	}

	redisCfg, err := config.RedisLoad(cfg.Environment)
	if err != nil {
		log.Fatalf("failed to connect to Redis: %w", err)
		return
	}
	redisClient, err := redis.NewRedisClient(redisCfg)
	if err != nil {
		log.Fatalf("Warning: Failed to connect to Redis: %v. Continuing without cache.", err)
		redisClient = nil
	}

	repos, err := repository.NewRepositories(cassandraDB, redisClient)
	if err != nil {
		log.Fatalf("Failed to initialize repositories: %v", err)
	}

	services, err := service.NewServices(&service.ServiceDependencies{Repos: repos})
	if err != nil {
		log.Fatalf("Failed to initialize service: %v", err)
	}

	middleware := middleware.NewMiddleware(services, repos)

	r := gin.Default()
	routes.SetupRoutes(&routes.Routes{
		Router:     r,
		Services:   services,
		Middleware: middleware,
	})
	r.Run(":8080")
}
