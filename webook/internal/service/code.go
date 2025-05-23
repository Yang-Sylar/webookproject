package service

import (
	"context"
	"fmt"
	"math/rand"
	"webook/internal/repository"
	"webook/internal/service/sms"
)

const codeTplid = "1877556"

var (
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
)

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, code string) (bool, error)
}

type MyCodeService struct {
	repo   *repository.CachedCodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo *repository.CachedCodeRepository, smsSvc sms.Service) *MyCodeService {
	return &MyCodeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

func (svc *MyCodeService) Send(ctx context.Context, biz string, phone string) error {

	// 码谁来生成，谁来管理？
	// 生成一个验证码
	code := svc.generateCode()

	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	err = svc.smsSvc.Send(ctx, codeTplid, []string{code}, phone)

	//if err != nil {
	//	// 意味着redis有验证码，但是用户没收到
	//	// 可以考虑重试，传入一个 retrySvc，也可以不管
	//	return err
	//}

	return err
}

func (svc *MyCodeService) Verify(ctx context.Context, biz string, phone string, code string) (bool, error) {

	return svc.repo.Verify(ctx, biz, phone, code)
}

func (svc *MyCodeService) generateCode() string {
	// 六位数，num 在 0，99999 之间，包含 0 和 99999
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
