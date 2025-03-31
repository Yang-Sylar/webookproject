package ratelimit

import (
	"context"
	"webook/internal/service/sms"
	"webook/pkg/limiter"
)

type RatelimitSMSService struct {
	svc     sms.Service
	limiter limiter.Limiter
}

func (s RatelimitSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		return err
	}

	if limited {
		return err
	}

	// 可以加新特性
	err = s.svc.Send(ctx, tpl, args, numbers...)
	// 可以加新特性
	if err != nil {
		return err
	}

	return nil
}
