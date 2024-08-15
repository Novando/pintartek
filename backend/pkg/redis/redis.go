package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/Novando/pintartek/pkg/logger"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	rdb    *redis.Client
	logger *logger.Logger
}

func Init(host string, port int, password string, logger *logger.Logger) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       0, // use default DB
	})

	status := rdb.Ping(context.Background())
	if status.Err() != nil {
		logger.Fatalf("Error connecting to redis: %s", status.Err())
	}

	logger.Info("Redis client initialized successfully.")

	return &Redis{
		rdb:    rdb,
		logger: logger,
	}
}

func (r *Redis) FlushAll() {
	r.rdb.FlushAll(context.Background())
}

func (r *Redis) Get(key string) (string, error) {
	val, err := r.rdb.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", fmt.Errorf("%s: %s", "Error getting value from redis", err)
	}
	return val, nil
}

func (r *Redis) Set(key string, value string, expiration time.Duration) error {
	_, err := r.rdb.Set(context.Background(), key, value, expiration).Result()
	if err != nil {
		return fmt.Errorf("%s: %s", "Error setting value in redis", err)
	}
	return nil
}

func (r *Redis) GetHash(key string, field string) (string, error) {
	val, err := r.rdb.HGet(context.Background(), key, field).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", fmt.Errorf("%s: %s", "Error getting value from redis", err)
	}
	return val, nil
}

func (r *Redis) SetHash(key string, field string, value string) error {
	_, err := r.rdb.HSet(context.Background(), key, field, value).Result()
	if err != nil {
		return fmt.Errorf("%s: %s", "Error setting value in redis", err)
	}
	return nil
}

func (r *Redis) SetHashTTL(key string, field string, value string, ttl time.Duration) error {
	_, err := r.rdb.HSet(context.Background(), key, field, value).Result()
	if err != nil {
		return fmt.Errorf("%s: %s", "Error setting value in redis", err)
	}

	err = r.rdb.Expire(context.Background(), field, ttl).Err()
	if err != nil {
		return fmt.Errorf("%s: %s", "Error setting TTL on hash", err)
	}

	return nil
}

func (r *Redis) Delete(key string) error {
	_, err := r.rdb.Del(context.Background(), key).Result()
	if err != nil {
		return fmt.Errorf("%s: %s", "Error deleting value from redis", err)
	}
	return nil
}

func (r *Redis) Close() {
	defer func() {
		err := r.rdb.Close()
		if err != nil {
			r.logger.Fatalf("%s: %s", "Error closing redis", err)
		}
	}()
}
