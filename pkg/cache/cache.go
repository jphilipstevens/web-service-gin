package cache

import (
	"context"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/jphilipstevens/web-service-gin/v2/pkg/apiErrors"
	"github.com/jphilipstevens/web-service-gin/v2/pkg/appTracer"
	"github.com/jphilipstevens/web-service-gin/v2/pkg/clientContext"
	"github.com/jphilipstevens/web-service-gin/v2/pkg/config"
	"github.com/jphilipstevens/web-service-gin/v2/pkg/datastore"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Cacher exposes Redis-style functionality while satisfying datastore.DataStore.
type Cacher interface {
	datastore.DataStore

	Get(serviceName string, ctx context.Context, key string) (val string, err error)
	Set(serviceName string, ctx context.Context, key string, value string, expiration time.Duration) error
}

type redisCache struct {
	Client    *redis.Client
	appTracer appTracer.AppTracer
}

var ErrCacheMiss = apiErrors.NewNotFoundError("")
var ErrCacheGeneric = apiErrors.NewGenericError("")

func NewCacher(cfg config.RedisClientConfig, appTracer appTracer.AppTracer) Cacher {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		panic(err)
	}
	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		panic(err)
	}

	return &redisCache{
		Client:    rdb,
		appTracer: appTracer,
	}
}

func (rc *redisCache) Get(serviceName string, ctx context.Context, key string) (string, error) {
	startTime := time.Now()
	ctx, span := rc.appTracer.CreateSpan(ctx, serviceName)
	defer span.End()

	val, err := rc.Client.Get(ctx, key).Result()
	hit := err == nil
	cacheCall := clientContext.CacheCall{
		ServiceTransaction: clientContext.ServiceTransaction{
			ServiceName: serviceName,
			SpanId:      span.SpanContext().TraceID().String(),
		},
		Action:       "get",
		ResponseTime: time.Since(startTime),
		Key:          key,
		Error:        err,
		Hit:          hit,
	}
	clientContext.AddCacheCall(ctx, cacheCall)
	dsCall := clientContext.DataStoreCall{
		ServiceTransaction: cacheCall.ServiceTransaction,
		StoreType:          "cache",
		Operation:          "get",
		Key:                key,
		Hit:                &hit,
		ResponseTime:       cacheCall.ResponseTime,
		Error:              err,
	}
	clientContext.AddDataStoreCall(ctx, dsCall)

	if err != nil {
		mappedErr := MapCacheError(&err)
		if err != redis.Nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		return "", mappedErr
	}

	span.SetStatus(codes.Ok, "")

	span.SetAttributes(attribute.String("cache.name", "redis"))
	span.SetAttributes(attribute.String("cache.action", "get"))
	span.SetAttributes(attribute.String("cache.key", key))
	span.SetAttributes(attribute.Int("cache.runeCount.", utf8.RuneCountInString(val)))

	return val, nil
}

func (rc *redisCache) Set(serviceName string, ctx context.Context, key string, value string, expiration time.Duration) error {
	startTime := time.Now()
	ctx, span := rc.appTracer.CreateSpan(ctx, serviceName)
	defer span.End()

	err := rc.Client.Set(ctx, key, value, expiration).Err()

	cacheCall := clientContext.CacheCall{
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
	clientContext.AddCacheCall(ctx, cacheCall)
	dsCall := clientContext.DataStoreCall{
		ServiceTransaction: cacheCall.ServiceTransaction,
		StoreType:          "cache",
		Operation:          "set",
		Key:                key,
		Hit:                nil,
		ResponseTime:       cacheCall.ResponseTime,
		Error:              err,
	}
	clientContext.AddDataStoreCall(ctx, dsCall)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return MapCacheError(&err)
	}
	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.String("cache.name", "redis"))
	span.SetAttributes(attribute.String("cache.action", "set"))
	span.SetAttributes(attribute.String("cache.key", key))
	span.SetAttributes(attribute.Int("cache.runeCount.", utf8.RuneCountInString(value)))
	span.SetAttributes(attribute.Int("cache.expirationSeconds", int(expiration.Seconds())))

	return MapCacheError(&err)
}

// ExecContext implements datastore.DataStore for write operations.
func (rc *redisCache) ExecContext(serviceName string, ctx context.Context, operation string, args ...any) (any, error) {
	if operation != "set" || len(args) < 3 {
		return nil, fmt.Errorf("unsupported cache exec: %s", operation)
	}
	key, _ := args[0].(string)
	value, _ := args[1].(string)
	exp, ok := args[2].(time.Duration)
	if !ok {
		return nil, fmt.Errorf("invalid expiration")
	}
	return nil, rc.Set(serviceName, ctx, key, value, exp)
}

// QueryContext implements datastore.DataStore for read operations.
func (rc *redisCache) QueryContext(serviceName string, ctx context.Context, operation string, args ...any) (any, error) {
	if operation != "get" || len(args) < 1 {
		return nil, fmt.Errorf("unsupported cache query: %s", operation)
	}
	key, _ := args[0].(string)
	return rc.Get(serviceName, ctx, key)
}

// Close shuts down the underlying redis client.
func (rc *redisCache) Close() {
	rc.Client.Close()
}

func MapCacheError(err *error) error {
	switch {
	case *err == redis.Nil:
		return nil
	case *err != nil:
		return ErrCacheGeneric
	default:
		return nil
	}
}
