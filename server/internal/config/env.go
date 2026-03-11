package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI                 string
	PORT                     int
	JWTSECRET                string
	UPSTASH_REDIS_REST_URL   string
	UPSTASH_REDIS_REST_TOKEN string
	UPSTASH_REDIS_URL        string
	REDIS_URL                string
	REDIS_ADDR               string
	REDIS_PASSWORD           string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	port, err := strconv.Atoi(strings.TrimSpace(os.Getenv("PORT")))

	if err != nil {
		return nil, fmt.Errorf("invalid PORT value: %v", err)
	}

	cfg := &Config{
		MongoURI:                 strings.TrimSpace(os.Getenv("MONGO_URI")),
		PORT:                     port,
		JWTSECRET:                strings.TrimSpace(os.Getenv("JWT_SECRET")),
		UPSTASH_REDIS_REST_URL:   strings.TrimSpace(os.Getenv("UPSTASH_REDIS_REST_URL")),
		UPSTASH_REDIS_REST_TOKEN: strings.TrimSpace(os.Getenv("UPSTASH_REDIS_REST_TOKEN")),
		UPSTASH_REDIS_URL:        strings.TrimSpace(os.Getenv("UPSTASH_REDIS_URL")),
		REDIS_URL:                strings.TrimSpace(os.Getenv("REDIS_URL")),
		REDIS_ADDR:               strings.TrimSpace(os.Getenv("REDIS_ADDR")),
		REDIS_PASSWORD:           strings.TrimSpace(os.Getenv("REDIS_PASSWORD")),
	}
	return cfg, nil
}
