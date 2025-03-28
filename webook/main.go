package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
	"webook/config"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"
	"webook/pkg/XtremeGin/middlewares/ratelimit"
)

func main() {
	db := initDB()
	redisClient := initRedis()

	server := initWebServer(redisClient)

	u := initUser(db, redisClient)
	u.RegisterRoutes(server)
	//server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "K8s部署测试")
	})
	server.Run(":8080")
}

func initWebServer(cmd redis.Cmdable) *gin.Engine {
	server := gin.Default()

	// 限流
	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: config.Config.Redis.Addr,
	//})
	// 限流
	server.Use(ratelimit.NewBuilder(cmd, time.Second, 100).Build())

	// Use 跨域问题，作用于全部路由
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},         // 允许请求的路由来源，* 为全部来源（一般不写）
		AllowMethods:     []string{"POST", "GET"},                   // 允许请求的方法，不写为全部方法
		AllowHeaders:     []string{"Content-Type", "Authorization"}, //
		AllowCredentials: true,                                      // 是否允许携带 cookie 之类的用户认证信息
		ExposeHeaders:    []string{"x-jwt-token"},                   // 响应头里携带的东西
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

	//// session基于redis实现
	//// 最大空闲连接、tcp（不太可能用到udp）、连接信息和密码、两个key
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
	//	// authentication key, encryption key（身份认证、数据加密）外加授权认证————安全三大概念
	//	[]byte("YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6T"),
	//	[]byte("e5Z7W4YbVcerrtjEA77eT5J6hShjjNTp"))
	//if err != nil {
	//	panic(err)
	//}

	store := memstore.NewStore(
		[]byte("YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6T"),
		[]byte("e5Z7W4YbVcerrtjEA77eT5J6hShjjNTp"))

	server.Use(sessions.Sessions("mysession", store)) // cookie的name和值
	//server.Use(
	//	middleware.
	//		NewLoginMiddlewareBuilder().
	//		IgnorePaths("/users/signup"). // 忽略路径
	//		IgnorePaths("/users/login").
	//		Build())

	server.Use(
		middleware.
			NewLoginJWTMiddlewareBuilder().
			IgnorePaths("/users/signup"). // 忽略路径
			IgnorePaths("/users/login").
			Build())

	return server
}

func initUser(db *gorm.DB, redis redis.Cmdable) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	ur := cache.NewUserCache(redis)
	repo := repository.NewUserRepository(ud, ur)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
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
