package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
}

func Load() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development" // 默認環境
	}
	// 加載對應環境的 .env 文件
	err := godotenv.Load(fmt.Sprintf(".env.%s", env))
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{
		Environment: env,
	}

	return config, nil
}
