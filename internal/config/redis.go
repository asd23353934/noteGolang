package config

import (
	"fmt"
	"os"
)

type RedisConfig struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	RedisPoolSize int
}

func RedisLoad(env string) (*RedisConfig, error) {
	config := &RedisConfig{
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       0,
		RedisPoolSize: 10,
	}

	if config.RedisHost == "" {
		return nil, fmt.Errorf("REDIS_HOST is not set")
	}
	// if config.RedisPassword == "" {
	// 	return nil, fmt.Errorf("REDIS_PASSWORD is not set")
	// }
	if config.RedisPort == "" {
		return nil, fmt.Errorf("REDIS_PORT is not set")
	}
	fmt.Println(config)
	return config, nil
}
