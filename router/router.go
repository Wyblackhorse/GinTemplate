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
	"github.com/wangyi/GinTemplate/controller"
	eeor "github.com/wangyi/GinTemplate/error"
)

func Setup() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)


	r := gin.New()



	r.Use(eeor.ErrHandler())
	r.NoMethod(eeor.HandleNotFound)
	r.NoRoute(eeor.HandleNotFound)
	r.Static("/static", "./static")
	//上传文件
	r.POST("/uploadFiles", controller.UploadFiles)
	//GetRecharge
	r.GET("/getRecharge", controller.GetRecharge)
	//GetUserApp
	r.GET("/getUserApp", controller.GetUserApp)
	//GetWalletRecord
	r.GET("/getWalletRecord", controller.GetWalletRecord)
	r.GET("/GetWithdraw", controller.GetWithdraw)

	r.Run(fmt.Sprintf(":%d", viper.GetInt("app.port")))

	return r
}
