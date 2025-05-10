package myjwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yzletter/go-kit/setx"
	"net/http"
)

type JwtServiceBuilder struct {
	IgnorePaths *setx.Set[string] // go-kit 工具库: github.com/yzletter/go-kit/
	JwtHandler  *JwtHandler
}

// NewJwtServiceBuilder 构造函数
func NewJwtServiceBuilder(jwtHandler *JwtHandler) *JwtServiceBuilder {
	return &JwtServiceBuilder{
		IgnorePaths: setx.NewSet[string](),
		JwtHandler:  jwtHandler,
	}
}

// AddIgnorePath 添加需要忽略不进行鉴权的路径
func (jb *JwtServiceBuilder) AddIgnorePath(path string) *JwtServiceBuilder {
	jb.IgnorePaths.Insert(path)
	return jb
}

// Build 返回用于 Gin 注册的 gin.HandlerFunc
func (jb *JwtServiceBuilder) Build() gin.HandlerFunc {
	var CheckRequestJWT func(ctx *gin.Context)

	CheckRequestJWT = func(ctx *gin.Context) {
		// 1. 判断当前路径是否需要忽略
		requestRoute := ctx.Request.URL.Path
		if jb.IgnorePaths.Exist(requestRoute) {
			return
		}
		// 2. 取出 tokenString
		tokenString := ExtractToken(ctx)
		// 3. 将 Token 解析成 Claims
		targetAccessClaims := &AccessClaims{}
		keyFunc := func(token *jwt.Token) (interface{}, error) {
			return jb.JwtHandler.AccessTokenKey, nil
		}
		token, err := jwt.ParseWithClaims(tokenString, targetAccessClaims, keyFunc)
		// 4. 判断是否携带 token
		if err != nil { // 未登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 5. 判断 Token 是否有效, 是否为空, Uid 是否为 0
		if !token.Valid || token == nil || targetAccessClaims.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 6. 判断 Agent 是否一致, 若不一致存在问题
		if targetAccessClaims.UserAgent != ctx.Request.UserAgent() {
			// 严重的安全问题，你是要监控
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 7. 判断 Token 是否废弃
		if ifDiscarded := jb.JwtHandler.CheckTokenDiscarded(ctx, targetAccessClaims.SSid); ifDiscarded {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 8. 将请求的用户数据存入上下文
		ctx.Set("claims", targetAccessClaims)

		// 9. 退出
		return
	}

	return CheckRequestJWT
}
