//go:build !k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "localhost:13306",
	},
	Redis: RedisConfig{
		Addr: "localhost:3309",
	},
}
