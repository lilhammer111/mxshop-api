package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/order-web/api/shop_cart"
	m "mxshop-api/order-web/middlewares"
)

func InitShopCartRouter(Router *gin.RouterGroup) {
	ShopCartRouter := Router.Group("shopcarts").Use(m.JWTAuth()).Use(m.Trace())
	zap.S().Info("配置用户相关url")
	{
		ShopCartRouter.GET("", shop_cart.List)
		ShopCartRouter.DELETE("/:id", shop_cart.Delete)
		ShopCartRouter.POST("", shop_cart.New)
		ShopCartRouter.PATCH("/:id", shop_cart.Update)
	}
}
