package cache

import (
	"context"
	"encoding/json"
	"github.com/GermanBogatov/user-service/internal/common/apperror"
	"github.com/GermanBogatov/user-service/internal/entity"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"time"
)

type ICache interface {
	Get(ctx context.Context, key string) (string, error)
	GetUser(ctx context.Context, key string) (entity.User, error)
	Delete(ctx context.Context, key string) error
	SetUser(ctx context.Context, key string, user entity.User) error
	SetRefreshToken(ctx context.Context, key, userID string) error
}

var _ ICache = &Cache{}

type Cache struct {
	client     *redis.Client
	userTTL    time.Duration
	refreshTTL time.Duration
}

func NewStorage(client *redis.Client, userTTL, refreshTTL int) ICache {
	return &Cache{
		client:     client,
		userTTL:    time.Duration(userTTL) * time.Second,
		refreshTTL: time.Duration(refreshTTL) * time.Second,
	}
}

// Get - получение данных из кеша по ключу.
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", apperror.ErrRedisNil
		}
		return "", err
	}

	return value, nil
}

// Delete - удаление записи из кэша по ключу
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// GetUser - получение пользователя из кэша
func (c *Cache) GetUser(ctx context.Context, key string) (entity.User, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return entity.User{}, apperror.ErrRedisNil
		}
		return entity.User{}, err
	}

	var user entity.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

// SetUser - добавление пользователя в кеш.
func (c *Cache) SetUser(ctx context.Context, key string, user entity.User) error {
	data, errJson := json.Marshal(user)
	if errJson != nil {
		return errJson
	}

	_, err := c.client.Set(ctx, key, string(data), c.userTTL).Result()
	if err != nil {
		return err
	}

	return nil
}

// SetRefreshToken - добавление рефреш токена в кэш.
func (c *Cache) SetRefreshToken(ctx context.Context, key, userID string) error {
	_, err := c.client.Set(ctx, key, userID, c.refreshTTL).Result()
	if err != nil {
		return err
	}

	return nil
}
