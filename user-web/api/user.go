package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	//claims, _ := c.Get("claims")
	//currentUser := claims.(*models.CustomClaims)
	//zap.S().Infof("访问用户的ID为%d", currentUser.ID)

	pn := c.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := c.DefaultQuery("pSize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{Pn: uint32(pnInt), PSize: uint32(pSizeInt)})
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
	// set clear false for conveniently testing
	verifyOk := store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, false)
	if !verifyOk {
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	// 调用接口，登陆逻辑
	rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	})
	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, gin.H{"mobile": "用户不存在"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"mobile": "登陆失败"})
			}
		}
	} else {
		// 用户存在，开始校验密码
		checkRsp, checkErr := global.UserSrvClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          passwordLoginForm.Password,
			EncryptedPassword: rsp.Password,
		})
		if checkErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "登陆失败"})
		} else {
			if checkRsp.Success {
				// 生成token
				helper.GenerateJWTWhenVerified(c, rsp.Id, rsp.Role, rsp.NickName)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "密码错误"})
			}
		}
	}
}

func Register(c *gin.Context) {
	// 表单验证
	registerForm := forms.RegisterForm{}
	err := c.ShouldBind(&registerForm)
	if err != nil {
		helper.HandleValidatorError(c, err)
		return
	}

	//验证码
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})

	value, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil || value != registerForm.Code {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码错误",
		})
		return
	}

	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		Password: registerForm.PassWord,
		Mobile:   registerForm.Mobile,
	})

	if err != nil {
		zap.S().Errorf("[Register] 查询 【用户列表】失败: %s", err.Error())
		helper.HandleGrpcErrorToHttp(err, c)
		return
	}

	helper.GenerateJWTWhenVerified(c, user.Id, user.Role, user.NickName)

}
