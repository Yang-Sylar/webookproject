package web

import "github.com/gin-gonic/gin"

// UserHandler 用于在上面定义所有跟user有关的路由
type UserHandler struct {
}

func (u *UserHandler) RegisterRoutesv1(ug *gin.RouterGroup) {
	ug.POST("/signup", u.Signup)  // 注册
	ug.POST("/login", u.Login)    // 登录
	ug.POST("/edit", u.Edit)      // 编辑
	ug.GET("/profile", u.Profile) // 个人信息
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	// 分组路由
	ug := server.Group("/users")

	ug.POST("/signup", u.Signup)  // 注册
	ug.POST("/login", u.Login)    // 登录
	ug.POST("/edit", u.Edit)      // 编辑
	ug.GET("/profile", u.Profile) // 个人信息
}

func (u *UserHandler) Signup(ctx *gin.Context) {

}

func (u *UserHandler) Login(ctx *gin.Context) {

}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {

}
