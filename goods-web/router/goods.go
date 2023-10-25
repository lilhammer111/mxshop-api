package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/goods-web/api/goods"
	m "mxshop-api/goods-web/middlewares"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods").Use(m.Trace()).Use(m.JWTAuth())
	//zap.S().Info("配置用户相关url")
	{
		GoodsRouter.GET("", goods.List)
		// The interface for adding new goods requires administrator privileges
		GoodsRouter.POST("", m.IsAdminAuth(), goods.New)
		GoodsRouter.GET("/:id", goods.Detail)
		GoodsRouter.DELETE("/:id", m.IsAdminAuth(), goods.Delete)
		GoodsRouter.GET("/:id/stocks", goods.Stock)
		GoodsRouter.PATCH("/:id", m.IsAdminAuth(), goods.UpdateStatus)
		GoodsRouter.PUT("/:id", m.IsAdminAuth(), goods.Update)
	}
}
