package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"telegram_bot/internal/datastruct"
	"time"
)

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	redisCache := new(RedisCache)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	redisCache.Client = client
	return redisCache
}

func (c RedisCache) Ping(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}

func (c RedisCache) Set(ctx context.Context, key string, value interface{}) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = c.Client.Set(ctx, key, raw, time.Hour).Err()
	return err
}

func (c RedisCache) GetContest(ctx context.Context, key string) (*datastruct.Contest, error) {
	value, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	fmt.Println(value)
	contest := new(datastruct.Contest)
	err = json.Unmarshal([]byte(value), &contest)
	if err != nil {
		return nil, err
	}
	return contest, nil
}

func (c RedisCache) Delete(ctx context.Context, key string) error {
	err := c.Client.Del(ctx, key).Err()
	return err
}

func (c RedisCache) GetAllContests(ctx context.Context) ([]*datastruct.Contest, error) {
	keys, err := c.Client.Keys(ctx, "contest-*").Result()
	if err != nil {
		return nil, err
	}
	contests := make([]*datastruct.Contest, len(keys))
	for i, key := range keys {
		contests[i], err = c.GetContest(ctx, key)
		if err != nil {
			return nil, err
		}
	}
	return contests, nil
}
