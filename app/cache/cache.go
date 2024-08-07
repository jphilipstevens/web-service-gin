package cache

import (
	"context"
	"example/web-service-gin/app/apiErrors"
	"example/web-service-gin/app/appTracer"
	"example/web-service-gin/app/clientContext"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type RedisClientConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Cacher interface {
	Get(serviceName string, ctx context.Context, key string) (val string, err error)
	Set(serviceName string, ctx context.Context, key string, value string, expiration time.Duration) error
}

type redisCache struct {
	Client *redis.Client
}

var ErrCacheMiss = apiErrors.NewNotFoundError("")
var ErrCacheGeneric = apiErrors.NewGenericError("")

func NewCacher(cfg RedisClientConfig) Cacher {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &redisCache{Client: rdb}
}

func (rc *redisCache) Get(serviceName string, ctx context.Context, key string) (string, error) {
	startTime := time.Now()
	ctx, span := appTracer.CreateDownstreamSpan(ctx, serviceName)
	defer span.End()

	val, err := rc.Client.Get(ctx, key).Result()
	currentContext := ctx.Value(clientContext.ClientContextKey).(*clientContext.ClientContext)
	newCacheCall := clientContext.CacheCall{
		ServiceTransaction: clientContext.ServiceTransaction{
			ServiceName: serviceName,
			SpanId:      span.SpanContext().TraceID().String(),
		},
		Action:       "get",
		ResponseTime: time.Since(startTime),
		Key:          key,
		Error:        err,
		Hit:          val == "",
	}
	currentContext.Cache = append(currentContext.Cache, newCacheCall)
	_ = context.WithValue(ctx, clientContext.ClientContextKey, currentContext)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", MapCacheError(&err)
	}
	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.String("cache.action", "get"))
	span.SetAttributes(attribute.String("cache.key", key))
	span.SetAttributes(attribute.String("cache.value", val))

	return val, nil
}

func (rc *redisCache) Set(serviceName string, ctx context.Context, key string, value string, expiration time.Duration) error {
	startTime := time.Now()
	ctx, span := appTracer.CreateDownstreamSpan(ctx, serviceName)
	defer span.End()

	err := rc.Client.Set(ctx, key, value, expiration).Err()

	currentContext := ctx.Value(clientContext.ClientContextKey).(*clientContext.ClientContext)
	newCacheCall := clientContext.CacheCall{
		ServiceTransaction: clientContext.ServiceTransaction{
			ServiceName: serviceName,
			SpanId:      span.SpanContext().TraceID().String(),
		},
		Action:       "set",
		ResponseTime: time.Since(startTime),
		Key:          key,
		Error:        err,
		Hit:          false,
	}
	currentContext.Cache = append(currentContext.Cache, newCacheCall)
	_ = context.WithValue(ctx, clientContext.ClientContextKey, currentContext)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return MapCacheError(&err)
	}
	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.String("cache.action", "set"))
	span.SetAttributes(attribute.String("cache.key", key))

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
