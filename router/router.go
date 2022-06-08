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
	r.Static("/static", "./static")
	GroupV2 := r.Group("/v2")
	{
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
		//资金归集
		GroupV2.GET("/collection", v2.Collection)
		//更新总余额
		GroupV2.GET("/getAllMoney", v2.GetAllMoney)
		//手动回调 HandBackStatus
		GroupV2.GET("/handBackStatus", v2.HandBackStatus)
		//登录 Login
		GroupV2.GET("/login", v2.Login)

		//测试接口

		GroupV2.POST("/backUrl", func(context *gin.Context) {


			type T struct {
				Code   int    `json:"Code"`
				Msg    string `json:"Msg"`
				Result struct {
					Data string `json:"Data"`
				} `json:"Result"`
			}

			var jsonOne  T
			err := context.BindJSON(&jsonOne)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			//fmt.Println(json)



			fmt.Println(jsonOne)

		})

	}
	r.GET("/getaddr", v2.Getaddr)
	r.POST("/getaddr", v2.Getaddr)

	r.Run(fmt.Sprintf(":%d", viper.GetInt("app.port")))
	return r
}

