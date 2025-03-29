package retryable

import (
	"context"
	"webook/internal/service/sms"
)

// 要小心并发问题
type Service struct {
	svc      sms.Service
	retrycnt int // 重试次数
}

func (s Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	err := s.Send(ctx, tpl, args, numbers...)
	for err != nil && s.retrycnt < 10 {
		s.retrycnt++
		err = s.Send(ctx, tpl, args, numbers...)
	}
	return err
}
