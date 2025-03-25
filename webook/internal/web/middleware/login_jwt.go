package middleware

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
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

		tokenStr := segs[1]
		token, err := jwt.Parse(tokenStr,
			func(token *jwt.Token) (interface{}, error) {
				return []byte("YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6T"), nil
			})

		if err != nil {
			// 没登录
			ctx.String(http.StatusUnauthorized, "未登录")
			return
		}

		if !token.Valid || token == nil {
			// 没登录
			ctx.String(http.StatusUnauthorized, "未登录")
			return
		}
	}
}
