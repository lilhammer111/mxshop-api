package api

import (
	"github.com/gin-gonic/gin"
	bc "github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"net/http"
)

var store = bc.DefaultMemStore

func GetCaptcha(c *gin.Context) {
	cp := bc.NewCaptcha(bc.DefaultDriverDigit, store)
	id, b64s, err := cp.Generate()
	if err != nil {
		zap.S().Errorf("生成验证码错误: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成验证码错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"picPath":   b64s,
	})
}
