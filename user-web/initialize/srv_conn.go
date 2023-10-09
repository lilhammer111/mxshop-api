package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"
)

func SrvConn() {
	// 服务发现
	consulInfo := global.ServerConfig.ConsulInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			consulInfo.Host,
			consulInfo.Port,
			global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		// 负载均衡
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalln("[SrvConn] 连接 【用户服务失败】")
	}

	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}

// 问题：后续用户服务下线了或者ip,端口变了，我就需要重新发现这个服务
// 答案：负载均衡可以解决

//func SrvConn() {
//	// 从注册中心获取到用户服务的信息
//	cfg := capi.DefaultConfig()
//	//cfg.Address = "127.0.0.1:8500"
//	cfg.Address = fmt.Sprintf("%s:%d",
//		global.ServerConfig.ConsulInfo.Host,
//		global.ServerConfig.ConsulInfo.Port)
//	cli, err := capi.NewClient(cfg)
//	if err != nil {
//		zap.S().Fatalln(err)
//	}
//	filter := fmt.Sprintf(`Service == "%s"`, global.ServerConfig.UserSrvInfo.Name)
//	services, err := cli.Agent().ServicesWithFilter(filter)
//	if err != nil {
//		panic(err)
//	}
//	userSrvHost := ""
//	userSrvPort := 0
//
//	for _, service := range services {
//		userSrvHost = service.Address
//		userSrvPort = service.Port
//		break
//	}
//
//	if userSrvHost == "" {
//		zap.S().Fatalln("[SrvConn] 连接 【用户服务失败】")
//		return
//	}
//	// 拨号连接用户grpc服务器
//	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost,
//		userSrvPort), grpc.WithInsecure())
//	if err != nil {
//		//zap.S().Errorw("[GetUserList] 连接 【用户服务】失败",
//		//	"msg", err.Error(),
//		//)
//		zap.S().Errorw("连接 【用户服务】失败",
//			"msg", err.Error(),
//		)
//	}
//	// 1. 后续用户服务下线 2. 改端口了 3. 改ip了
//	userSrvClient := proto.NewUserClient(userConn)
//	global.UserSrvClient = userSrvClient
//}
