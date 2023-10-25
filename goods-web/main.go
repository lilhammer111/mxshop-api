package main

import (
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/initialize"
	"mxshop-api/goods-web/utils"
	"mxshop-api/goods-web/utils/register/consul"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	initialize.Logger()

	initialize.Config()

	r := initialize.Routers()

	if err := initialize.Translation("zh"); err != nil {
		zap.S().Fatal("初始化翻译器错误")
	}

	initialize.SrvConn()

	initialize.TrafficLimitRules()

	viper.AutomaticEnv()
	// 如果是本地开发环境则固定端口
	debug := viper.GetBool("MXSHOP_DEBUG")
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	// goods-web service register
	registryClient := consul.NewRegistryClient(
		global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port,
	)
	serviceID, err := uuid.GenerateUUID()
	zap.S().Infoln("serviceID is ", serviceID)

	if err != nil {
		zap.S().Errorf("failed to generate serviceID: %s", err.Error())
	}

	err = registryClient.Register(
		global.ServerConfig.Host,
		global.ServerConfig.Name,
		serviceID,
		global.ServerConfig.Port,
		global.ServerConfig.Tags,
	)
	if err != nil {
		zap.S().Panic("fail to register goods-web service", err.Error())
	}

	// run service
	go func() {
		zap.S().Debugf("在:%d上成功启动服务器", global.ServerConfig.Port)
		err = r.Run(fmt.Sprintf(":%d", global.ServerConfig.Port))
		if err != nil {
			zap.S().Panic("启动失败：", err.Error())
		}
	}()

	// receive quit signal
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	err = registryClient.Deregister(serviceID)
	if err != nil {
		zap.S().Panic("fail to deregister", err.Error())
	} else {
		zap.S().Infoln("succeed to deregister")
	}
}
