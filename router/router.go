/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package router

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	v2 "github.com/wangyi/GinTemplate/controller/v2"
	eeor "github.com/wangyi/GinTemplate/error"
	"github.com/wangyi/GinTemplate/tools"
	"github.com/wangyi/GinTemplate/util"
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

			var jsonOne T
			err := context.BindJSON(&jsonOne)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			//对date 进行解密
			decodeString, err1 := base64.StdEncoding.DecodeString(jsonOne.Result.Data)
			if err1 != nil {
				tools.ReturnError101(context, "err:"+err1.Error())
				return
			}

			one, err := util.RsaDecryptForEveryOne(decodeString)
			if err != nil {
				return
			}

			type CreatePrepaidPhoneOrdersData struct {
				PlatformOrder    string
				RechargeAddress  string
				Username         string
				AccountOrders    float64 //订单充值金额
				AccountPractical float64 //  实际充值的金额
				RechargeType     string
				BackUrl          string
			}

			var jsonData CreatePrepaidPhoneOrdersData
			err3 := json.Unmarshal(one, &jsonData)
			if err3 != nil {
				return
			}

			fmt.Println(jsonData)

		})
		GroupV2.GET("/TestBackUrl", func(context *gin.Context) {
			type Create struct {
				PlatformOrder    string
				RechargeAddress  string
				Username         string
				AccountOrders    float64 //订单充值金额
				AccountPractical float64 //  实际充值的金额
				RechargeType     string
				BackUrl          string
			}

			var tt Create
			tt.PlatformOrder = "1111111111111111111111111111"
			tt.RechargeAddress = "TW2HWaLWy9pwiRN4yLju6YKW3aQ6Fw8888"
			tt.Username = "wine"
			tt.AccountOrders = 100.00
			tt.AccountPractical = 10.00
			tt.RechargeType = "USDT"
			data, err := json.Marshal(tt)
			if err != nil {
				fmt.Println(err.Error())
			}
			data, _ = util.RsaEncryptForEveryOne(data)

			util.BackUrlToPay("http://127.0.0.1:7777/v2/backUrl", base64.StdEncoding.EncodeToString(data))

			tools.ReturnError200(context, "调佣成功")

		})
		//唤醒订单   PullUpTheOrder
	//	GroupV2.POST("/PullUpTheOrder", v2.PullUpTheOrder)

	}
	r.GET("/getaddr", v2.Getaddr)
	r.POST("/getaddr", v2.Getaddr)

	r.Run(fmt.Sprintf(":%d", viper.GetInt("app.port")))
	return r
}
