package myjwt

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	CheckTokenDiscarded(ctx *gin.Context, SSid string) bool
	ClearToken(ctx *gin.Context) error
	SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error
	SetAccessToken(ctx *gin.Context, uid int64, ssid string) error
	SetLoginToken(ctx *gin.Context, uid int64) error
	RefreshAccessToken(ctx *gin.Context)
}

// RefreshClaims 长 Token 声明
type RefreshClaims struct {
	Uid  int64
	SSid string
	jwt.RegisteredClaims
}

// AccessClaims 短 Token 声明
type AccessClaims struct {
	Uid       int64
	SSid      string
	UserAgent string
	jwt.RegisteredClaims
}

var (
	ErrSetAccessToken  = errors.New("XtremeGin: JWT 设置 AccessToken 错误")
	ErrSetRefreshToken = errors.New("XtremeGin: JWT 设置 RefreshToken 错误")
	ErrRedisSetSSid    = errors.New("XtremeGin: Redis 设置 SSid 错误")
)
