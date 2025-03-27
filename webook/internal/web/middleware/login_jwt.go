package middleware

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
	"webook/internal/web"
)

// LoginJWTMiddlewareBuilder JWT登录
type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {

		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			// 没登录
			ctx.String(http.StatusUnauthorized, "未登录")
			return
		}

		segs := strings.SplitN(tokenHeader, " ", 2)
		if len(segs) != 2 {
			// 没登录
			ctx.String(http.StatusUnauthorized, "未登录")
			return
		}

		claims := &web.UserClaims{}
		tokenStr := segs[1]
		// ParseWithClaims 里一定要传claims指针
		token, err := jwt.ParseWithClaims(
			tokenStr,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte("YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6T"), nil
			})

		if err != nil {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !token.Valid || token == nil || claims.Uid == 0 {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims.UserAgent != ctx.Request.UserAgent() {
			// 严重的安全问题
			// 你是要监控
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims.ExpiresAt.Sub(time.Now()) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour))
			tokenStr, err = token.SignedString([]byte("YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6T"))
			if err != nil {
				// 记录日志
				// 续约失败
			}
			ctx.Header("x-jwt-token", tokenStr)
		}
		ctx.Set("claims", claims)
	}
}
