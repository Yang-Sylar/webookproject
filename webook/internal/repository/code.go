package repository

import (
	"context"
	"webook/internal/repository/cache"
)

var (
	ErrCodeSendTooMany        = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
)

type CodeRepository interface {
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
	Store(ctx context.Context, biz string, phone string, code string) error
}

type CachedCodeRepository struct {
	cache *cache.RedisCodeCache
}

func NewCachedCodeRepository(cache *cache.RedisCodeCache) *CachedCodeRepository {
	return &CachedCodeRepository{
		cache: cache,
	}
}

func (repo *CachedCodeRepository) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, inputCode)
}

func (repo *CachedCodeRepository) Store(ctx context.Context, biz string, phone string, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}
