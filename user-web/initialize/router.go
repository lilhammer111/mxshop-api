package initialize

import (
	"github.com/gin-gonic/gin"
	m "mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/router"
	"net/http"
)

func Routers() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	//r := gin.New()
	//r.Use(gin.Recovery())
	r := gin.Default()
	// 配置CORS
	r.Use(m.Cors())
	// health check api
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	ApiGroup := r.Group("/u/v1")
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)
	return r
}
