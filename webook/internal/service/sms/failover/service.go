package failover

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"webook/internal/service/sms"
)

type FailoverSMSService struct {
	svcs []sms.Service
	idx  uint64
}

func NewFailoverSMSService(svcs []sms.Service) sms.Service {
	return &FailoverSMSService{
		svcs: svcs,
	}
}

// 简单轮询2
func (f *FailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[int(i%length)]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			return err
		}

	}

	return errors.New("全部服务商都失败了")
}

// 简单轮询
func (f *FailoverSMSService) Sendv1(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tpl, args, numbers...)
		if err == nil {
			return nil
		}
		log.Println(err)
	}

	return errors.New("全部服务商都失败了")
}
