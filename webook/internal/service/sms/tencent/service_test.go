package tencent

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
	"testing"
)

// 这个需要手动跑，也就是你需要在本地搞好这些环境变量
func TestSender(t *testing.T) {
	secretId, ok := os.LookupEnv("")
	if !ok {
		t.Fatal()
	}
	secretKey, ok := os.LookupEnv("")

	c, err := sms.NewClient(common.NewCredential(secretId, secretKey),
		"",
		profile.NewClientProfile())
	if err != nil {
		t.Fatal(err)
	}

	s := NewService(c, "", "")

	testCases := []struct {
		name    string
		tplId   string
		params  []string
		numbers []string
		wantErr error
	}{
		{
			name:   "",
			tplId:  "",
			params: []string{""},
			// 改成你的手机号码
			numbers: []string{""},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			er := s.Send(context.Background(), tc.tplId, tc.params, tc.numbers...)
			assert.Equal(t, tc.wantErr, er)
		})
	}
}
