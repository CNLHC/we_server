package main

import (
	"net/http"
	"os"
	"we_server/pkg/wx"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	viper.SetDefault("WX_APPID", "AID")
	viper.SetDefault("WX_APPSECRET", "ASECRET")
	viper.SetDefault("WX_APPTOKEN", "TOKEN")
	var redis_default = make(map[string]string)
	redis_default["Host"] = "REDIS_HOST"
	redis_default["Port"] = "REDIS_PORT"

	viper.SetDefault("redis", redis_default)

}

func GetAccessToken(c *gin.Context) {
	account := wx.GetAccount()
	if token, err := account.GetAccessToken(); err != nil {
		logrus.Errorf("can not get token: %s", err)
		c.JSON(200, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"data": token})
	}

}

type JSSDKReq struct {
	Uri string `form:"uri"`
}

func GetJSSDK(c *gin.Context) {
	var req JSSDKReq
	var err error
	if err = c.ShouldBind(&req); err == nil {
		logrus.Infof("get jsssdk rea: %+v", req)
		account := wx.GetAccount()
		js := account.GetJs()
		jscfg, err := js.GetConfig(req.Uri)
		if err != nil {
			goto err_handle
		}
		c.JSON(http.StatusOK, gin.H{"data": jscfg})
		return
	}
err_handle:
	c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

}

func main() {
	r := gin.Default()
	sublog := zerolog.New(os.Stdout)
	r.Use(logger.SetLogger(logger.Config{
		UTC:    true,
		Logger: &sublog,
	}))

	r.Use(cors.Default())
	g := r.Group("/api/v1")
	{
		g.GET("token", GetAccessToken)
		g.GET("jssdk/cfg", GetJSSDK)
	}

	r.Run(":2134")

}
