package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"webook/internal/domain"
)

var ErrUserNotFound = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
}

type RedisUserCache struct {
	// 传单机Redis可以，传cluster的Redis也可以
	client     redis.Cmdable
	expiration time.Duration
}

// A 用到 B，B 一定是接口
// A 用到 B，B 一定是 A 的字段
// A 用到 B，A 绝对不初始化 B，而是注入
func NewRedisUserCache(client redis.Cmdable) *RedisUserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

//type RedisCache interface {
//	Get(ctx context.Context, id int64)
//	Set(ctx context.Context, id int64, expiration time.Duration)
//}

// 如果没有数据，返回一个特定的err
func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {

	key := cache.key(id)
	val, err := cache.client.Get(ctx, key).Bytes()

	if err != nil {
		return domain.User{}, ErrUserNotFound
	}

	var u domain.User

	// 反序列化
	err = json.Unmarshal(val, &u)
	return u, err
}

func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	// 序列化
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
