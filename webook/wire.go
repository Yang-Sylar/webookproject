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
		// User DAO
		dao.NewUserDAO,
		// Cache
		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,
		// Repository
		repository.NewCachedCodeRepository,
		repository.NewUserRepository,
		// Service
		service.NewUserService,
		service.NewCodeService,

		// Web
		web.NewUserHandler,

		ioc.InitSMSService,
		ioc.InitMiddlewares,
		ioc.InitGin,
	)
	return new(gin.Engine)
}
