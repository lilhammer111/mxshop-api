package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mxshop-api/order-web/global"
	"mxshop-api/order-web/proto"
)

func SrvConn() {
	// 服务发现
	consulInfo := global.ServerConfig.ConsulInfo
	orderConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			consulInfo.Host,
			consulInfo.Port,
			global.ServerConfig.OrderSrvInfo.Name),
		grpc.WithInsecure(),
		// 负载均衡
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalln("[SrvConn] dialogue【order-srv failed】")
	}
	global.OrderSrvClient = proto.NewOrderClient(orderConn)

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
		zap.S().Fatalln("[SrvConn] dialogue【goods-srv failed】")
	}
	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)

	invConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			consulInfo.Host,
			consulInfo.Port,
			global.ServerConfig.InventorySrvInfo.Name),
		grpc.WithInsecure(),
		// 负载均衡
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalln("[SrvConn] dialogue【inventory-srv failed】")
	}
	global.InventorySrvClient = proto.NewInventoryClient(invConn)
}
