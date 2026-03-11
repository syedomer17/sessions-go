package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI  string
	PORT      int
	JWTSECRET string
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
		MongoURI: strings.TrimSpace(os.Getenv("MONGO_URI")),
		PORT:      port,
		JWTSECRET: strings.TrimSpace(os.Getenv("JWT_SECRET")),
	}
	return cfg, nil
}