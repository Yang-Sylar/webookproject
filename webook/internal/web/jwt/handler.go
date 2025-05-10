package myjwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
	"time"
)

type JwtHandler struct {
	AccessTokenKey  []byte
	RefreshTokenKey []byte
	RedisClient     redis.Cmdable
}

func NewJwtHandler(redisClient redis.Cmdable) *JwtHandler {
	return &JwtHandler{
		AccessTokenKey:  []byte("YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6T"),
		RefreshTokenKey: []byte("YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6a"),
		RedisClient:     redisClient,
	}
}

// CheckTokenDiscarded 判断 Token 是否被废弃
func (jh *JwtHandler) CheckTokenDiscarded(ctx *gin.Context, SSid string) bool {
	cnt, _ := jh.RedisClient.Exists(ctx, fmt.Sprintf("users:ssid:%s", SSid)).Result()
	if cnt > 0 {
		return true
	}
	return false
}

// SetLoginToken 设置登录 Token, 包括长 token 和短 token
func (jh *JwtHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := jh.SetRefreshToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = jh.SetAccessToken(ctx, uid, ssid)
	return err
}

// SetRefreshToken 设置长 token
func (jh *JwtHandler) SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	// 1. 携带信息的声明
	myClaims := RefreshClaims{
		Uid:  uid,
		SSid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "webook",                                               // 签名人
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 过期时间
		},
	}
	// 2. 生成 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, myClaims)
	// 3. 对 Token 进行加密
	tokenString, err := token.SignedString(jh.RefreshTokenKey) // 用长 token 秘钥进行加密
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return ErrSetRefreshToken
	}
	// 4. 将 token 放入上下文
	ctx.Header("x-refresh-token", tokenString)
	return nil
}

// SetAccessToken 设置短 token
func (jh *JwtHandler) SetAccessToken(ctx *gin.Context, uid int64, ssid string) error {
	// 1. 携带信息的声明
	myClaims := AccessClaims{
		Uid:       uid,
		SSid:      ssid,
		UserAgent: ctx.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "webook",                                           // 签名人
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 过期时间
		},
	}
	// 2. 生成 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, myClaims)
	// 3. 对 Token 进行加密
	tokenString, err := token.SignedString(jh.AccessTokenKey) // 用短 token 秘钥进行加密
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return ErrSetAccessToken
	}
	// 4. 将 token 放入上下文
	ctx.Header("x-jwt-token", tokenString)

	return nil
}

// ClearToken 将 token 废弃
func (jh *JwtHandler) ClearToken(ctx *gin.Context) error {
	// 1. 设置前端请求头的长短 token 为非法值
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	// 2. 获取当前请求的 SSid
	myClaims := ctx.MustGet("claims").(*AccessClaims)
	// 3. 在 Redis 中记录当前 SSid 废弃
	err := jh.RedisClient.Set(ctx, fmt.Sprintf("users:ssid:%s", myClaims.SSid), "", time.Hour*24*7).Err()
	if err != nil {
		return ErrRedisSetSSid
	}
	return nil
}

// RefreshAccessToken 若长 token 未过期, 则刷新短 token
func (jh *JwtHandler) RefreshAccessToken(ctx *gin.Context) {
	// 1. 取出 tokenString
	tokenString := ExtractToken(ctx)
	// 2. 用 refreshTokenKey 解析到声明中
	targetClaims := &RefreshClaims{}
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return jh.RefreshTokenKey, nil
	}
	token, err := jwt.ParseWithClaims(tokenString, targetClaims, keyFunc)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 3. 判断 token 是否废弃
	if ifDiscarded := jh.CheckTokenDiscarded(ctx, targetClaims.SSid); ifDiscarded {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 4. 判断 token 是否正确
	if !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 5. 设置新的 access token
	err = jh.SetAccessToken(ctx, targetClaims.Uid, targetClaims.SSid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 6. 返回
	ctx.String(http.StatusOK, "刷新 accessToken 成功")
}

// ExtractToken 从上下文取出 tokenString
func ExtractToken(ctx *gin.Context) string {
	headerString := ctx.GetHeader("Authorization")
	headerStringSeg := strings.SplitN(headerString, " ", 2)
	if len(headerStringSeg) != 2 {
		return ""
	}
	return headerStringSeg[1]
}
