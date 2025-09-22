package repository

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrNotFound - стандартная ошибка, когда ключ в кэше не найден.
var ErrNotFound = errors.New("key not found in cache")

// CacheRedis реализует CacheRepository с использованием Redis.
type CacheRedis struct {
	client *redis.Client
}

// NewCacheRedis создает новый экземпляр репозитория кэша.
func NewCacheRedis(client *redis.Client) *CacheRedis {
	return &CacheRedis{client: client}
}

// Set сохраняет значение в кэше с временем жизни (TTL).
func (r *CacheRedis) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Get извлекает значение из кэша. Возвращает ErrNotFound, если ключ не существует.
func (r *CacheRedis) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrNotFound
		}
		return "", err
	}
	return val, nil
}

// Delete удаляет ключ из кэша.
func (r *CacheRedis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
