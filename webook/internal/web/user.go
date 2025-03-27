package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
	"webook/internal/domain"
	"webook/internal/service"
)

// UserHandler 用于在上面定义所有跟user有关的路由
type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	// 正则表达式
	const (
		emailRegexPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
		passwordRegexPattern = `^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,16}$`
	)
	// 预编译
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	// 返回指针
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	// 分组路由
	ug := server.Group("/users")

	ug.POST("/signup", u.Signup) // 注册
	ug.POST("/login", u.Login)   // 登录

	ug.POST("/edit", u.Edit)      // 编辑
	ug.GET("/profile", u.Profile) // 个人信息
}

func (u *UserHandler) Signup(ctx *gin.Context) {
	// 函数内部定义结构体不让其他人用
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var req SignUpReq
	// Bind 方法会根据 Content-type 来解析到你的 req 里面
	// 解析错了就会写回一个 400 的错误，传指针
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 邮箱校验
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}

	// 密码校验
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码格式错误")
		return
	}

	// 调用一下service的方法
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱重复，请换一个邮箱")
		return
	} else if err != nil {
		ctx.String(http.StatusOK, "服务器异常，注册失败")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrInvaildUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码错误")
		return
	} else if err != nil {
		ctx.String(http.StatusOK, "服务器异常，注册失败")
		return
	}

	//// 在这里成功登陆了
	//sess := sessions.Default(ctx)
	//// 设置session的值
	//sess.Set("userId", user.Id)
	//sess.Options(sessions.Options{
	//	MaxAge:   30,    // 消亡时间
	//	Secure:   false, // 要求https协议
	//	HttpOnly: false, // 只允许http
	//})
	//sess.Save()
	claims := UserClaims{
		Uid:       user.Id,
		UserAgent: ctx.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), // 过期时间
			Issuer:    "webook",                                      // 签发人
		},
	}

	// 用JWT实现登录态
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6T"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Save()
	ctx.String(http.StatusOK, "退出成功")
	return
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}

	var req EditReq
	// 1. 获取参数
	//{nickname: "yzleter", birthday: "2025-03-27", aboutMe: "golang backend engineer"}
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}

	// 2. 校验参数
	if len(req.Nickname) > 20 {
		ctx.String(http.StatusOK, "昵称过长")
		return
	}

	TimeLayout := "2006-01-02"
	bdate, err := time.Parse(TimeLayout, req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "日期错误")
		return
	}

	if len(req.AboutMe) > 200 {
		ctx.String(http.StatusOK, "简介过长")
		return
	}

	// 拿 userid
	c, ok := ctx.Get("claims")
	if !ok {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	claims, ok := c.(*UserClaims)
	fmt.Println(claims.Uid)
	fmt.Println(req)

	// 3. 处理业务
	err = u.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		Id:       claims.Uid,
		Nickname: req.Nickname,
		Birthday: bdate,
		AboutMe:  req.AboutMe,
	})

	if err != nil {
		ctx.String(http.StatusInternalServerError, "修改失败")
		return
	}

	// 4. 退出
	ctx.String(http.StatusOK, "修改成功")
	return
}

func (u *UserHandler) Profile(ctx *gin.Context) {

	c, ok := ctx.Get("claims")
	if !ok {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}

	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}

	uback, err := u.svc.GetProfile(ctx, claims.Uid)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return

	}

	type User struct {
		Nickname string `json:"Nickname"`
		Email    string `json:"Email"`
		AboutMe  string `json:"AboutMe"`
		Birthday string `json:"Birthday"`
	}

	ctx.JSON(http.StatusOK, User{
		Nickname: uback.Nickname,
		Email:    uback.Email,
		AboutMe:  uback.AboutMe,
		Birthday: uback.Birthday.Format(time.DateOnly),
	})
}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	c, ok := ctx.Get("claims")
	// 可以断定必然有claims
	if !ok {
		// 可以考虑监控住这里
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	claims, ok := c.(*UserClaims) // 断言
	if !ok {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	println(claims.Uid)
	// Profile 其他代码
	ctx.String(http.StatusOK, "这是你的Profile")
}

type UserClaims struct {
	Uid       int64 // 声明自己要放进去 token 里面的数据
	UserAgent string
	jwt.RegisteredClaims
}
