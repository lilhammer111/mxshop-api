package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
	"mxshop-api/user-web/utils"
	"mxshop-api/user-web/utils/register/consul"
	customValidator "mxshop-api/user-web/validator"
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

	viper.AutomaticEnv()
	// 如果是本地开发环境则固定端口
	debug := viper.GetBool("MXSHOP_DEBUG")
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	// register validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", customValidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	// user-web service register
	registryClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceID, _ := uuid.GenerateUUID()
	err := registryClient.Register(global.ServerConfig.Host, global.ServerConfig.Name, serviceID, global.ServerConfig.Port, global.ServerConfig.Tags)
	if err != nil {
		zap.S().Panic("fail to register user-web service", err.Error())
	}

	go func() {
		zap.S().Debugf("启动服务器，端口为： %d", global.ServerConfig.Port)

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
