package repository

import (
	"context"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
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
	u, err := r.dao.FindById(ctx, id)

	if err != nil {
		return domain.User{}, err
	}
	
	return domain.User{
		Email:    u.Email,
		Nickname: u.Nickname,
		Birthday: time.UnixMilli(u.Birthday),
		AboutMe:  u.AboutMe,
	}, err

	//先从 cache 找
	//再从数据库找
	//找到写回 cache
}
