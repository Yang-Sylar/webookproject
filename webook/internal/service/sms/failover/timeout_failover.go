package failover

import (
	"context"
	"sync/atomic"
	"webook/internal/service/sms"
)

type TimeoutFailoverSMSService struct {

	// 服务商
	svcs []sms.Service
	idx  int32
	cnt  int32

	// 阈值：连续响应超时超过这个数字就切换
	threshold int32
}

func (t TimeoutFailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	length := len(t.svcs)

	if cnt > t.threshold {
		// 切换
		newIdx := (idx + 1) % int32(length)
		// CompareAndSwapInt32(地址，旧值，新值）进行交换
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			// 我成功挪了一位，即指向了新服务商，累计访问次数清零
			atomic.StoreInt32(&t.cnt, 0)
		}
		// else 就是出现并发了
		idx = atomic.LoadInt32(&t.idx)
	}

	svc := t.svcs[idx]
	err := svc.Send(ctx, tpl, args, numbers...)
	switch err {
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
	default:
		// 不知道什么错误
		// 你可以考虑 切
		return err
	}

	return err
}
