package web

import "github.com/gin-gonic/gin"

// UserHandler 用于在上面定义所有跟user有关的路由
type UserHandler struct {
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	// 注册
	server.POST("/users/signup", u.Signup)
	// 登录
	server.POST("/users/login", u.Signup)
	// 编辑
	server.POST("/users/edit", u.Edit)
	// 个人信息
	server.GET("/users/profile", u.Profile)
}

func (u *UserHandler) Signup(ctx *gin.Context) {

}
func (u *UserHandler) Login(ctx *gin.Context) {

}
func (u *UserHandler) Edit(ctx *gin.Context) {

}
func (u *UserHandler) Profile(ctx *gin.Context) {

}
