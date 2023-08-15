package main

import (
	"fmt"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
)

func main() {
	//port := 8021
	// initialize logger
	initialize.Logger()
	// initialize config
	initialize.Config()
	// initialize routers
	r := initialize.Routers()
	// initialize translation
	if err := initialize.Translation("zh"); err != nil {
		zap.S().Fatal("初始化翻译器错误")
	}

	zap.S().Debugf("启动服务器，端口为： %d", global.ServerConfig.Port)

	err := r.Run(fmt.Sprintf(":%d", global.ServerConfig.Port))
	if err != nil {
		zap.S().Panic("启动失败：", err.Error())
	}
}
