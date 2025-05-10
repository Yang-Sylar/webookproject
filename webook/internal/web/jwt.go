package web

// import (
//
//	"github.com/gin-gonic/gin"
//	"github.com/golang-jwt/jwt/v5"
//	"github.com/google/uuid"
//	"net/http"
//	"strings"
//	"time"
//
// )
//
//	type jwtHandler struct {
//		// accessToken key
//		atKey []byte
//		// refreshToken key
//		rtKey []byte
//	}
//
//	func newJWTHandler() jwtHandler {
//		return jwtHandler{
//			atKey: []byte("YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6T"),
//			rtKey: []byte("YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6a"),
//		}
//	}
//
//	type RefreshClaims struct {
//		Uid  int64
//		Ssid string
//		jwt.RegisteredClaims
//	}
//func (jh jwtHandler) setLoginToken(ctx *gin.Context, uid int64) error {
//	ssid := uuid.New().String()
//	err := jh.setJWTToken(ctx, uid, ssid)
//	if err != nil {
//		return err
//	}
//	err = jh.setRefreshToken(ctx, uid, ssid)
//	return err
//}

//
//func (jh jwtHandler) setJWTToken(ctx *gin.Context, uid int64, ssid string) error {
//
//	claims := UserClaims{
//		Uid:       uid,
//		UserAgent: ctx.Request.UserAgent(),
//		Ssid:      ssid,
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), // 过期时间
//			Issuer:    "webook",                                      // 签发人
//		},
//	}
//
//	// 用JWT实现登录态
//	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
//	tokenStr, err := token.SignedString(jh.atKey)
//	if err != nil {
//		ctx.String(http.StatusInternalServerError, "系统错误")
//		return err
//	}
//	ctx.Header("x-jwt-token", tokenStr)
//	return nil
//}
//
//func (jh jwtHandler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
//
//	claims := RefreshClaims{
//		Uid:  uid,
//		Ssid: ssid,
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 10)), // 过期时间
//			Issuer:    "webook",                                           // 签发人
//		},
//	}
//
//	// 用JWT实现登录态
//	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
//	tokenStr, err := token.SignedString(jh.rtKey)
//	if err != nil {
//		ctx.String(http.StatusInternalServerError, "系统错误")
//		return err
//	}
//	ctx.Header("x-refresh-token", tokenStr)
//	return nil
//}
//
//func ExtractToken(ctx *gin.Context) string {
//	tokenHeader := ctx.GetHeader("Authorization")
//	//fmt.Println(tokenHeader)
//	segs := strings.SplitN(tokenHeader, " ", 2)
//	//fmt.Println(segs)
//	if len(segs) != 2 {
//		return ""
//	}
//	//fmt.Println(segs[1])
//	return segs[1]
//}
