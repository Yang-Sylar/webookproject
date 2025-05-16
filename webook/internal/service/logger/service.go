package logger

import (
	"context"
	"go.uber.org/zap"
	"webook/internal/service/sms"
)

type Service struct {
	svc sms.Service
}

func (s *Service) Send(ctx context.Context, biz string, args []string, numbers ...string) {
	zap.L().Debug("发送信息", zap.String("biz", biz))
	err := s.svc.Send(ctx, biz, args, numbers...)
	if err != nil {
		zap.L().Debug("发送信息", zap.Error(err))
	}
}
