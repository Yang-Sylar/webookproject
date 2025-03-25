package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {

		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		sess := sessions.Default(ctx)

		id := sess.Get("userId")
		if id == nil {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		sess.Options(sessions.Options{
			MaxAge: 30, // 消亡时间
		})
		now := time.Now()
		updateTime := sess.Get("update_time") // 获取上次刷新时间
		sess.Set("userId", id)
		if updateTime == nil {
			// 说明还没刷新过
			sess.Set("update_time", now)
			if err := sess.Save(); err != nil {
				panic(err)
			}
			return
		}

		updateTimeVal, _ := updateTime.(time.Time)

		if now.Sub(updateTimeVal) > time.Second*10 {
			// 该刷新了
			sess.Set("update_time", now)
			if err := sess.Save(); err != nil {
				panic(err)
			}
			return
		}
	}
}
