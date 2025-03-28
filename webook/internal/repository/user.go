package repository

import (
	"context"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (r *UserRepository) UpdateNonZeroFields(ctx context.Context, u domain.User) error {
	err := r.dao.UpdateById(ctx, dao.User{
		Id:       u.Id,
		Nickname: u.Nickname,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
	})
	if err != nil {
		return err
	}
	return err
	// 在这里操作缓存
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})

	// 在这里操作缓存
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, err
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {

	// 先从缓存找
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, err
	}

	// 缓存没这个数据
	if err == cache.ErrUserNotFound {
		// 去数据库里加载
		// 考虑 Redis 可能崩了，大量访问直接把数据库也崩了
	}

	ur, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	u = domain.User{
		Id:       ur.Id,
		Email:    ur.Email,
		Nickname: ur.Nickname,
		Birthday: time.UnixMilli(ur.Birthday),
		AboutMe:  ur.AboutMe,
	}

	//go func() {
	//	err = r.cache.Set(ctx, u)
	//	if err != nil {
	//		// 做好监控
	//	}
	//}()
	err = r.cache.Set(ctx, u)
	if err != nil {
		// 做好监控
	}
	return u, err

}
