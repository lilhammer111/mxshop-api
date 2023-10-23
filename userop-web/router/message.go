package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/userop-web/api/message"
	m "mxshop-api/userop-web/middlewares"
)

func InitMessageRouter(Router *gin.RouterGroup) {
	MessageRouter := Router.Group("message").Use(m.JWTAuth()).Use(m.Trace())
	{
		MessageRouter.GET("", message.List) // 轮播图列表页
		MessageRouter.POST("", message.New) //新建轮播图
	}
}
