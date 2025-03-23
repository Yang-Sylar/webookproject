package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserHandler 用于在上面定义所有跟user有关的路由
type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler() *UserHandler {
	const (
		emailRegexPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
		passwordRegexPattern = `^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,16}$`
	)
	// 预编译
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	// 返回指针
	return &UserHandler{
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

//func (u *UserHandler) RegisterRoutesv1(ug *gin.RouterGroup) {
//	ug.POST("/signup", u.Signup)  // 注册
//	ug.POST("/login", u.Login)    // 登录
//	ug.POST("/edit", u.Edit)      // 编辑
//	ug.GET("/profile", u.Profile) // 个人信息
//}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	// 分组路由
	ug := server.Group("/users")

	ug.POST("/signup", u.Signup)  // 注册
	ug.POST("/login", u.Login)    // 登录
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

	ctx.String(http.StatusOK, "注册成功")
	// 这边就是数据库操作
}

func (u *UserHandler) Login(ctx *gin.Context) {

}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {

}
