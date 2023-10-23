package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/order-web/api/order"
	m "mxshop-api/order-web/middlewares"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	OrderRouter := Router.Group("orders").Use(m.JWTAuth()).Use(m.Trace())
	zap.S().Info("配置用户相关url")
	{
		OrderRouter.GET("", order.List)
		OrderRouter.POST("", order.New)
		OrderRouter.GET("/:id", order.Detail)
	}
}
