package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/goods-web/api/goods"
	m "mxshop-api/goods-web/middlewares"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	zap.S().Info("配置用户相关url")
	{
		GoodsRouter.GET("", goods.List)
		// The interface for adding new goods requires administrator privileges
		GoodsRouter.POST("", m.JWTAuth(), m.IsAdminAuth(), goods.New)
	}
}
