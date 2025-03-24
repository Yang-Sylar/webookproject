package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"
)

func main() {
	db := initDB()
	server := initWebServer()

	u := initUser(db)
	u.RegisterRoutes(server)

	server.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	// Use 作用于全部路由
	server.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     []string{"http://localhost:3000"},         // 允许请求的路由来源，* 为全部来源（一般不写）
				AllowMethods:     []string{"POST", "GET"},                   // 允许请求的方法，不写为全部方法
				AllowHeaders:     []string{"Content-Type", "Authorization"}, //
				AllowCredentials: true,                                      // 是否允许携带 cookie 之类的用户认证信息
				//ExposeHeaders:    []string{"x-jwt-token"},	// 响应头里携带的东西
				// 判断来源的函数
				AllowOriginFunc: func(origin string) bool {
					if strings.Contains(origin, "http://localhost") {
						// 开发环境
						return true
					}
					return strings.Contains(origin, "youcompany.com")
				},
				MaxAge: 12 * time.Hour,
			}))

	// 步骤1
	// session
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mysession", store)) // cookie的name和值
	
	server.Use(
		middleware.
			NewLoginMiddlewareBuilder().
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login").
			Build())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		// 我只会在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 换言之，一旦初始化出错，就别启动了
		panic(err)
	}

	// 建表
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	return db
}
