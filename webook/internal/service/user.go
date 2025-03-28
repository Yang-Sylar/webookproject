package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"webook/internal/domain"
	"webook/internal/repository"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvaildUserOrPassword = errors.New("账号或密码不对")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 考虑加密放在哪里的问题
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	// 然后就是存起来
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, u domain.User) (domain.User, error) {
	// 先找用户
	r, err := svc.repo.FindByEmail(ctx, u.Email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvaildUserOrPassword
	}

	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(r.Password), []byte(u.Password))
	if err != nil {
		// DEBUG
		return domain.User{}, ErrInvaildUserOrPassword
	}
	return r, err
}

func (svc *UserService) UpdateNonSensitiveInfo(ctx context.Context, u domain.User) error {
	return svc.repo.UpdateNonZeroFields(ctx, u)
}

func (svc *UserService) GetProfile(ctx context.Context, id int64) (domain.User, error) {
	u, err := svc.repo.FindById(ctx, id)
	return u, err
}
