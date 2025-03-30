package service

import (
	"context"
	"fmt"
	"math/rand"
	"webook/internal/repository"
	"webook/internal/service/sms"
)

const codeTplid = "1877556"

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo *repository.CodeRepository, smsSvc sms.Service) *CodeService {
	return &CodeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

func (svc *CodeService) Send(ctx context.Context,
	// 区别业务场景
	biz string,
	phone string) error {

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

func (svc *CodeService) Verify(ctx context.Context,
	// 区别业务场景
	biz string,
	code string, phone string) (bool, error) {

	return svc.repo.Verify(ctx, biz, phone, code)
}

func (svc *CodeService) generateCode() string {
	// 六位数，num 在 0，99999 之间，包含 0 和 99999
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
