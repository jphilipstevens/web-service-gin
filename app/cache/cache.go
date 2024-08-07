package cache

import (
	"context"
	"example/web-service-gin/app/apiErrors"
	"example/web-service-gin/app/appTracer"
	"example/web-service-gin/app/clientContext"
	"example/web-service-gin/config"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type Cacher interface {
	Get(serviceName string, ctx context.Context, key string) (val string, err error)
	Set(serviceName string, ctx context.Context, key string, value string, expiration time.Duration) error
}

type redisCache struct {
	Client *redis.Client
}

var ErrCacheMiss = apiErrors.NewNotFoundError("")
var ErrCacheGeneric = apiErrors.NewGenericError("")

func NewCacher(cfg config.RedisClientConfig) Cacher {
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
	newCacheCall := clientContext.CacheCall{
		ServiceTransaction: clientContext.ServiceTransaction{
			ServiceName: serviceName,
			SpanId:      span.SpanContext().TraceID().String(),
		},
		Action:       "get",
		ResponseTime: time.Since(startTime),
		Key:          key,
		Error:        err,
		Hit:          val != "",
	}
	clientContext.AddCacheCall(ctx, newCacheCall)
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
	clientContext.AddCacheCall(ctx, newCacheCall)

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
