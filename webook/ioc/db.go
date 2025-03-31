package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"webook/config"
	"webook/internal/repository/dao"
)

func InitDB() *gorm.DB {
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
