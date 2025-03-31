//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/ioc"
)

func initWebServer() *gin.Engine {

	wire.Build(
		// 最基础的
		ioc.InitRedis,
		ioc.InitDB,

		dao.NewUserDAO,

		cache.NewCodeCache,
		cache.NewUserCache,

		repository.NewCachedCodeRepository,
		repository.NewUserRepository,

		service.NewUserService,
		service.NewCodeService,

		ioc.InitSMSService,
		web.NewUserHandler,
		ioc.InitMiddlewares,
		ioc.InitGin,
	)
	return new(gin.Engine)
}
