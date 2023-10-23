package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/user-web/api"
	m "mxshop-api/user-web/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user").Use(m.Trace())
	zap.S().Info("配置用户相关url")
	{
		UserRouter.GET("list", m.JWTAuth(), m.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("pwd_login", api.LoginByPWD)
		UserRouter.POST("register", api.Register)
	}
}
