package cache

import (
	"github.com/redis/go-redis/v9"
	"time"
)

type UserCache struct {
	// 传单机Redis可以，传cluster的Redis也可以
	client     redis.Cmdable
	expiration time.Duration
}

// A 用到 B，B 一定是接口
// A 用到 B，B 一定是 A 的字段
// A 用到 B，A 绝对不初始化 B，而是注入
//func NewUserCache(client redis.Cmdable) *UserCache {
//	return &UserCache{
//		client:     client,
//		expiration: time.Minute * 15,
//	}
//}
//
//func (u *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
//	u.client.Get(id)
//}
//
//func (u *UserCache) Set(ctx context.Context, id int64) (domain.User, error) {
//	u.client.Set(ctx, id)
//}
