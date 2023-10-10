package initialize

import (
	"github.com/gin-gonic/gin"
	m "mxshop-api/goods-web/middlewares"
	"mxshop-api/goods-web/router"
	"net/http"
)

func Routers() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	//r := gin.New()
	//r.Use(gin.Recovery())
	r := gin.Default()
	// 配置CORS
	r.Use(m.Cors())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	ApiGroup := r.Group("/g/v1")
	router.InitGoodsRouter(ApiGroup)
	return r
}
