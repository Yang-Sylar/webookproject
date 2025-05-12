package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {
	//initViper()
	//initViperRemote()
	initViperV1()
	server := initWebServer()
	server.Run(":8080")
}

func initViper() {
	viper.SetDefault("db.mysql.dsn", "root:root@tcp(localhost:3306)/webook")
	viper.SetConfigFile("config/dev.yaml")

	// 另一种写法
	//viper.SetConfigName("dev")      // 配置文件名, 不包含 .yaml 等后缀
	//viper.SetConfigType("yaml")     // 指定格式
	//viper.AddConfigPath("./config") // 当前工作目录下的config子目录

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

// 远程配置中心
func initViperRemote() {
	viper.SetConfigType("yaml")
	err := viper.AddRemoteProvider("etcd3",
		// 通过 webook 和其他使用 etcd 的区别出来
		"127.0.0.1:12379", "D:/Git/webook")
	if err != nil {
		panic(err)
	}

	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

// 获取启动参数
func initViperV1() {
	cfile := pflag.String("config", "config/config.yaml", "配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)

	viper.WatchConfig() // 实时监听配置变更
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 只能告诉你改了, 具体改了什么不清楚, 得重新读
		fmt.Println(in.Name, in.Op)
	})
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}
