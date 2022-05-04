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
	GroupV2.POST("/getPayInformation", v2.GetPayInformationBack)
	GroupV2.GET("/getPayInformation", v2.GetPayInformation)
	//资金归集
	//GroupV2.GET("/getPayInformation", v2.GetPayInformation)

	//充值订单管理
	GroupV2.POST("/createPrepaidPhoneOrders", v2.CreatePrepaidPhoneOrders)
	GroupV2.GET("/getPrepaidPhoneOrders", v2.GetPrepaidPhoneOrders)

	//GetReceiveAddress 地址管理
	GroupV2.GET("/getReceiveAddress", v2.GetReceiveAddress)

	//每日统计 DailyStatistics
	GroupV2.GET("/getDailyStatistics", v2.GetDailyStatistics)

	//测试接口
	r.GET("/getaddr", v2.Getaddr)

	r.Static("/static", "./static")
	r.Run(fmt.Sprintf(":%d", viper.GetInt("app.port")))
	return r
}
