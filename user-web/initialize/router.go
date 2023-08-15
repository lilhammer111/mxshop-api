package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/router"
)

func Routers() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	//r := gin.New()
	//r.Use(gin.Recovery())
	r := gin.Default()
	ApiGroup := r.Group("/u/v1")
	router.InitUserRouter(ApiGroup)
	return r
}
