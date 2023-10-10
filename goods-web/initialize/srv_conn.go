package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/proto"
)

func SrvConn() {
	// 服务发现
	consulInfo := global.ServerConfig.ConsulInfo
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			consulInfo.Host,
			consulInfo.Port,
			global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		// 负载均衡
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalln("[SrvConn] 连接 【用户服务失败】")
	}

	goodsSrvClient := proto.NewGoodsClient(goodsConn)
	global.GoodsSrvClient = goodsSrvClient
}
