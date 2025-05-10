package middleware

//
//import (
//	"encoding/gob"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/golang-jwt/jwt/v5"
//	"github.com/redis/go-redis/v9"
//	"net/http"
//	"time"
//	"webook/internal/web"
//)
//
//// LoginJWTMiddlewareBuilder JWT登录
//type LoginJWTMiddlewareBuilder struct {
//	paths []string
//	cmd   redis.Cmdable
//}
//
//func NewLoginJWTMiddlewareBuilder(cmd redis.Cmdable) *LoginJWTMiddlewareBuilder {
//	return &LoginJWTMiddlewareBuilder{
//		cmd: cmd,
//	}
//}
//
//func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
//	gob.Register(time.Now())
//	return func(ctx *gin.Context) {
//
//		for _, path := range l.paths {
//			if ctx.Request.URL.Path == path {
//				return
//			}
//		}
//
//		tokenStr := web.ExtractToken(ctx)
//		claims := &web.UserClaims{}
//		// ParseWithClaims 里一定要传claims指针
//		token, err := jwt.ParseWithClaims(
//			tokenStr,
//			claims,
//			func(token *jwt.Token) (interface{}, error) {
//				return []byte("YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6T"), nil
//			})
//
//		if err != nil {
//			// 没登录
//			ctx.AbortWithStatus(http.StatusUnauthorized)
//			return
//		}
//
//		if !token.Valid || token == nil || claims.Uid == 0 {
//			// 没登录
//			ctx.AbortWithStatus(http.StatusUnauthorized)
//			return
//		}
//
//		if claims.UserAgent != ctx.Request.UserAgent() {
//			// 严重的安全问题，你是要监控
//			ctx.AbortWithStatus(http.StatusUnauthorized)
//			return
//		}
//
//		ctx.Set("claims", claims)
//
//		count, err := l.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", claims.Ssid)).Result()
//		if err != nil || count > 0 {
//			// 要么 redis 有问题, 要么当前登录已经退出
//			ctx.AbortWithStatus(http.StatusUnauthorized)
//			return
//		}
//	}
//}
//
//func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
//	l.paths = append(l.paths, path)
//	return l
//}
