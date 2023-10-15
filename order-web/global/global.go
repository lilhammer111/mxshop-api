package global

import (
	ut "github.com/go-playground/universal-translator"
	"mxshop-api/order-web/config"
	"mxshop-api/order-web/proto"
)

var (
	ServerConfig   = &config.ServerConfig{}
	Trans          ut.Translator
	OrderSrvClient proto.OrderClient
	GoodsSrvClient proto.GoodsClient
	NacosConfig    = &config.NacosConfig{}
)
