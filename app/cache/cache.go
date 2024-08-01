package cache

import (
	"context"
	"example/web-service-gin/app/apiErrors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClientConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Cacher interface {
	Get(ctx context.Context, key string) (val string, err error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}

type redisClient struct {
	client *redis.Client
}

var ErrCacheMiss = apiErrors.NewNotFoundError("")
var ErrCacheGeneric = apiErrors.NewGenericError("")

func NewCacher(cfg RedisClientConfig) Cacher {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &redisClient{client: rdb}
}

func (rc *redisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		// TODO: Handle cache miss error
		return "", MapCacheError(&err)
	}
	return val, nil
}

func (rc *redisClient) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := rc.client.Set(ctx, key, value, expiration).Err()
	return MapCacheError(&err)
}

func MapCacheError(err *error) error {
	switch {
	case *err == redis.Nil:
		return ErrCacheMiss
	case *err != nil:
		return ErrCacheGeneric
	default:
		return nil
	}
}
