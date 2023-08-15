package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mxshop-api/user-web/api/helper"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/global/response"
	"mxshop-api/user-web/proto"
	"net/http"
	"strconv"
	"time"
)

func GetUserList(c *gin.Context) {
	//zap.S().Debug("获取用户列表页")

	//ip := "127.0.0.1"
	//port := 50051
	// 拨号连接用户grpc服务器
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务】失败",
			"msg", err.Error(),
		)
	}

	// 调用接口
	userSrvClient := proto.NewUserClient(userConn)
	pn := c.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := c.DefaultQuery("pSize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{Pn: uint32(pnInt), PSize: uint32(pSizeInt)})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表】失败")
		helper.HandleGrpcErrorToHttp(err, c)
		return
	}
	res := make([]response.UserResponse, 0)
	for _, value := range rsp.Data {
		//data := make(map[string]any)
		//data["id"] = value.Id
		//data["name"] = value.NickName
		//data["birthday"] = value.Birthday
		//data["gender"] = value.Gender
		//data["mobile"] = value.Mobile
		//res = append(res, data)

		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Birthday: time.Unix(int64(value.Birthday), 0).Format("2006-01-02"),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}

		res = append(res, user)
	}
	c.JSON(http.StatusOK, res)
}

func LoginByPWD(c *gin.Context) {
	//表单验证
	passwordLoginForm := forms.PWDLoginForm{}
	err := c.ShouldBind(&passwordLoginForm)
	if err != nil {
		helper.HandleValidatorError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "sign in succeeded",
	})
}
