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
	v1 "github.com/wangyi/GinTemplate/controller/v1"
	eeor "github.com/wangyi/GinTemplate/error"
)

func Setup() *gin.Engine {

	r := gin.New()
	r.Use(eeor.ErrHandler())
	r.NoMethod(eeor.HandleNotFound)
	r.NoRoute(eeor.HandleNotFound)
	r.Static("/static", "./static")

	v1Group := r.Group("/v1")

	{
		v1Group.GET("/setWebName", v1.SetWebName)
		v1Group.GET("/checkWebName", v1.CheckWebName)
	}

	r.Run(fmt.Sprintf(":%d", viper.GetInt("app.port")))
	return r
}
