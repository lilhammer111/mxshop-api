package helper

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/models"
	"net/http"
	"strings"
	"time"
)

const (
	// 30天的秒数
	JWTExpirationInterval = 60 * 60 * 24 * 30
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	// 将grpc的code转换成http的状态码
	e, ok := status.FromError(err)
	if ok {
		switch e.Code() {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"msg": e.Message()})
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "内部错误"})
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		case codes.Unavailable:
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "服务不可用"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "其他错误" + e.Message(), "code": e.Code()})
		}
		//return
	}
}

func RemoveTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for k, v := range fields {
		rsp[k[strings.Index(k, ".")+1:]] = v
	}
	return rsp
}

func HandleValidatorError(c *gin.Context, err error) {
	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": RemoveTopStruct(validationErr.Translate(global.Trans)),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	}
}

// GenerateJWTWhenVerified generate a jwt when the password was verified or users registered, and it
// response user's id, nickname, token and the expired time.
func GenerateJWTWhenVerified(c *gin.Context, id, role int32, nickname string) {
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(id),
		NickName:    nickname,
		AuthorityId: uint(role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + JWTExpirationInterval,
			Issuer:    "imooc",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "生成token失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         uint(id),
		"nick-name":  nickname,
		"token":      token,
		"expired-at": claims.ExpiresAt * 1000, // 毫秒级别
		//"msg":        "登陆成功",
	})
}
