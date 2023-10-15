package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/order-web/api/order"
	m "mxshop-api/order-web/middlewares"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	OrderRouter := Router.Group("orders")
	zap.S().Info("配置用户相关url")
	{
		OrderRouter.GET("", m.JWTAuth(), m.IsAdminAuth(), order.List)
		OrderRouter.POST("", m.JWTAuth(), order.New)
		OrderRouter.GET("/:id", m.JWTAuth(), order.Detail)
	}
}
