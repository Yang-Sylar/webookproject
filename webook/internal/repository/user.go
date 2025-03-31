package repository

import (
	"context"
	"database/sql"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository interface {
	UpdateNonZeroFields(ctx context.Context, u domain.User)
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
}

type CachedUserRepository struct {
	dao   *dao.GORMUserDAO
	cache *cache.RedisUserCache
}

func NewUserRepository(dao *dao.GORMUserDAO, cache *cache.RedisUserCache) *CachedUserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (r *CachedUserRepository) UpdateNonZeroFields(ctx context.Context, u domain.User) error {
	err := r.dao.UpdateById(ctx, dao.User{
		Id:       u.Id,
		Nickname: u.Nickname,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
	})
	if err != nil {
		return err
	}

	// 在这里更新缓存
	err = r.cache.Set(ctx, u)
	if err != nil {
		// 做好监控
	}

	return err

}

func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))
	// 在这里操作缓存
}

func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), err
}

func (r *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), err
}

func (r *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {

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

	u = r.entityToDomain(ur)
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

// domain 转 dao
func (r *CachedUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Ctime:    u.Ctime.UnixMilli(),
	}
}

// dao 转 domain
func (r *CachedUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
		Ctime:    time.UnixMilli(u.Ctime),
		Nickname: u.Nickname,
		Birthday: time.UnixMilli(u.Birthday),
		AboutMe:  u.AboutMe,
	}
}
