//go:build k8s

// 编译k8s的时候编译该文件
package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(webook-mysql:3309)/webook",
	},

	Redis: RedisConfig{
		Addr: "webook-redis:10379",
	},
}
