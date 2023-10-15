package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/order-web/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

// 从配置文件中读取出配置

func Config() {
	debug := GetEnvInfo("MXSHOP_DEBUG")
	zap.S().Infoln("env MXSHOPO_DEBUG is ", debug)
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("order-web/%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("order-web/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()
	// 文件路径如何设置？
	v.SetConfigFile(configFileName)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = v.Unmarshal(global.NacosConfig)
	if err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息: %v", global.NacosConfig)

	// 动态监控变化
	//v.WatchConfig()
	//v.OnConfigChange(func(e fsnotify.Event) {
	//	zap.S().Infof("配置文件产生变化: %s", e.Name)
	//	_ = v.ReadInConfig()
	//	_ = v.Unmarshal(global.ServerConfig)
	//	zap.S().Infof("配置信息: %v", global.ServerConfig)
	//})

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   uint64(global.NacosConfig.Port),
		},
	}
	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// 创建动态配置客户端的另一种方式 (推荐)
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group})

	if err != nil {
		panic(err)
	}

	//fmt.Println(content)

	err = json.Unmarshal([]byte(content), &global.ServerConfig)

	if err != nil {
		zap.S().Fatalf("读取nacos配置失败： %s\n", err)
	}

	//fmt.Printf("%+v", &global.ServerConfig)

}
