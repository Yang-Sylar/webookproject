package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type jwtHandler struct {
	// accessToken key
	atKey []byte
	// refreshToken key
	rtKey []byte
}

func newJWTHandler() jwtHandler {
	return jwtHandler{
		atKey: []byte("YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6T"),
		rtKey: []byte("YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6a"),
	}
}

type RefreshClaims struct {
	Uid int64
	jwt.RegisteredClaims
}

func (jh jwtHandler) setJWTToken(ctx *gin.Context, uid int64) error {

	claims := UserClaims{
		Uid:       uid,
		UserAgent: ctx.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), // 过期时间
			Issuer:    "webook",                                      // 签发人
		},
	}

	// 用JWT实现登录态
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(jh.atKey)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (jh jwtHandler) setRefreshToken(ctx *gin.Context, uid int64) error {

	claims := RefreshClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), // 过期时间
			Issuer:    "webook",                                      // 签发人
		},
	}

	// 用JWT实现登录态
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(jh.rtKey)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

func ExtractToken(ctx *gin.Context) string {
	tokenHeader := ctx.GetHeader("Authorization")
	segs := strings.SplitN(tokenHeader, " ", 2)
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}
