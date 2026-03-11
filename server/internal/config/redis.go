package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
	"strings"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

func init() {
	if err := ConnectRedis(); err != nil {
		panic(err)
	}
}

func ConnectRedis() error {

	cfg, err := Load()

	if err != nil {
		return err
	}

	redisURL := firstNonEmpty(cfg.REDIS_URL, cfg.UPSTASH_REDIS_URL, cfg.UPSTASH_REDIS_REST_URL)
	password := firstNonEmpty(cfg.REDIS_PASSWORD, cfg.UPSTASH_REDIS_REST_TOKEN)

	if redisURL == "" {
		return fmt.Errorf("missing redis endpoint: set REDIS_URL/UPSTASH_REDIS_URL or REDIS_ADDR")
	}

	var opt *redis.Options
	if strings.HasPrefix(redisURL, "redis://") || strings.HasPrefix(redisURL, "rediss://") {
		opt, err = redis.ParseURL(redisURL)
		if err != nil {
			return fmt.Errorf("invalid redis url %q: %w", redisURL, err)
		}
		if opt.Password == "" {
			opt.Password = password
		}
	} else if strings.HasPrefix(redisURL, "http://") || strings.HasPrefix(redisURL, "https://") {
		opt, err = optionsFromUpstashREST(redisURL, password)
		if err != nil {
			return err
		}
	} else {
		addr, addrErr := normalizeRedisAddr(firstNonEmpty(cfg.REDIS_ADDR, redisURL))
		if addrErr != nil {
			return addrErr
		}
		opt = &redis.Options{Addr: addr, Password: password}
	}

	RedisClient = redis.NewClient(opt)

	_, err = RedisClient.Ping(ctx).Result()

	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func normalizeRedisAddr(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", fmt.Errorf("missing redis address")
	}

	if strings.Contains(value, "://") {
		parsed, err := url.Parse(value)
		if err != nil {
			return "", fmt.Errorf("invalid redis address %q: %w", value, err)
		}

		switch parsed.Scheme {
		case "tcp", "redis", "rediss":
			addr := parsed.Host
			if addr == "" {
				addr = strings.TrimPrefix(parsed.Path, "/")
			}
			if !strings.Contains(addr, ":") {
				return "", fmt.Errorf("redis address must include host:port, got %q", value)
			}
			return addr, nil
		default:
			return "", fmt.Errorf("unsupported redis scheme %q", parsed.Scheme)
		}
	}

	if !strings.Contains(value, ":") {
		return "", fmt.Errorf("redis address must include host:port, got %q", value)
	}

	return value, nil
}

func optionsFromUpstashREST(restURL string, password string) (*redis.Options, error) {
	parsed, err := url.Parse(strings.TrimSpace(restURL))
	if err != nil {
		return nil, fmt.Errorf("invalid upstash rest url %q: %w", restURL, err)
	}

	if parsed.Host == "" {
		return nil, fmt.Errorf("invalid upstash rest url %q: missing host", restURL)
	}

	addr := parsed.Host
	if !strings.Contains(addr, ":") {
		addr = addr + ":6379"
	}

	opt := &redis.Options{
		Addr:     addr,
		Username: "default",
		Password: password,
	}

	if parsed.Scheme == "https" {
		opt.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	return opt, nil
}
