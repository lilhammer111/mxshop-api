package helper

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/user-web/global"
	"net/http"
	"strings"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	// 将grpc的code转换成http的状态码
	if err != nil {
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
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "用户服务不可用"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "其他错误" + e.Message(), "code": e.Code()})
			}
			//return
		}
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
