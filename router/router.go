/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	v2 "github.com/wangyi/GinTemplate/controller/v2"
	eeor "github.com/wangyi/GinTemplate/error"
)

func Setup() *gin.Engine {

	r := gin.New()
	r.Use(eeor.ErrHandler())
	r.NoMethod(eeor.HandleNotFound)
	r.NoRoute(eeor.HandleNotFound)

	GroupV2 := r.Group("/v2")
	GroupV2.POST("/getPayInformation",v2.GetPayInformationBack)

	r.Static("/static", "./static")
	r.Run(fmt.Sprintf(":%d", viper.GetInt("app.port")))
	return r
}
