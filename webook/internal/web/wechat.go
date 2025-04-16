package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webook/internal/service/oauth2/wechat"
)

type OAuth2WechatHandler struct {
	svc wechat.Service
}

func NewOAuth2WechatHandler(svc wechat.Service) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc: svc,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	url, err := h.svc.AuthURL(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "构造扫码登录url失败",
		})
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})

}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	ctx.String(http.StatusOK, "你过来啦")
}
