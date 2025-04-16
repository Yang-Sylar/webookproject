package wechat

import (
	"context"
	"fmt"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/url"
)

var redirectURI = url.PathEscape("www.baidu.com")

type Service interface {
	AuthURL(ctx context.Context) (string, error)
}

type service struct {
	appId string
}

func NewService(appId string) Service {
	return &service{
		appId: appId,
	}
}

func (s *service) AuthURL(ctx context.Context) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redire"
	state := uuid.New()
	return fmt.Sprintf(urlPattern, s.appId, redirectURI, state), nil
}
