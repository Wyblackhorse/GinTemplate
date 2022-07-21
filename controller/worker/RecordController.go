/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package worker

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"go.uber.org/zap"
	"strconv"
)

//获取账单
func GetRecord(c *gin.Context) {
	action := c.Query("action")
	who, _ := c.Get("who")
	whoMap := who.(map[string]interface{})
	//获取账单
	if action == "GET" {
		kinds, _ := strconv.Atoi(c.Query("kinds"))
		record := make([]model.Record, 0)
		mysql.DB.Where("kinds=? and worker_id=?", kinds, whoMap["ID"]).Order("created desc").Find(&record)
		ReturnSuccessData(c, record, "success")
		return
	}
	//提现
	if action == "withdraw" {
		//获取提现金额
		money, _ := strconv.ParseFloat(c.Query("money"), 64)
		//获取系统配置 最低的 提现金额
		config := model.Config{}
		err := mysql.DB.Where("id=?", 1).First(&config).Error
		if err != nil {
			ReturnErr101(c, "system is  fail ")
			return
		}
		if config.LowWithdrawal > money {
			ReturnErr101(c, "The minimum withdrawal amount is "+strconv.FormatFloat(config.LowWithdrawal, 'f', 2, 64))
			return
		}
		//金额校验完毕
		worker := model.WorkerBalance{ID: int(whoMap["ID"].(uint)), Kinds: 2, AddBalance: money}
		result, err1 := worker.AddBalanceFuc(mysql.DB)
		if result == false {
			ReturnErr101(c, err1.Error())
			return
		}
		d := model.DailyStatistics{TodayRechargeMoney: money, TodayRechargePeople: 1, WorkerId: int(whoMap["ID"].(uint))}
		d.SetEverydayData(mysql.DB)

		ReturnSuccess(c, "Withdrawal successful, waiting for administrator review")
		return
	}
	//充值
	if action == "recharge" {

		kinds, _ := strconv.Atoi(c.Query("kinds"))
		//d := model.DailyStatistics{TodayRechargeMoney: money,TodayRechargePeople: 1,WorkerId: int(whoMap["ID"].(uint))}
		//d.SetEverydayData(mysql.DB)
		//获取充值金额
		money, _ := strconv.ParseFloat(c.Query("money"), 64)
		if money <= 0 {
			ReturnErr101(c, "money  is  fail")
			return
		}
		//判断 充值平台类型
		if kinds == 1 {
			var pp PlatformOrderParameterize
			pp.Username = viper.GetString("PlatformOrderParameterize.Username")
			pp.AccountOrders = money
			pp.RechargeType = "USDT"
			pp.BackUrl = viper.GetString("PlatformOrderParameterize.BackUrl")
			r := model.Record{Kinds: 1, WorkerId: int(whoMap["ID"].(uint)),Money: money}
			pp.PlatformOrder = r.CratedNewRechargeOrder(mysql.DB)
			if pp.PlatformOrder == "" {
				ReturnErr101(c, "sorry  system is  wrong")
				return
			}
			//加密
			data, err := json.Marshal(pp)
			if err != nil {
				ReturnErr101(c, "sorry  system is  wrong1")
				return
			}
			data, err = tools.RsaEncrypt(data)
			if err != nil {
				ReturnErr101(c, "sorry  system is  wrong2")
				return
			}
			type postData struct {
				Data string `json:"data"`
			}
			//数据base64
			dataForbes64 := base64.StdEncoding.EncodeToString(data)
			var pd postData
			pd.Data = dataForbes64
			dataT, err := json.Marshal(pd)
			if err != nil {
				ReturnErr101(c, "sorry  system is  wrong3")
				return
			}
			rd, err1 := tools.PostHttp(viper.GetString("PlatformOrderParameterize.PrepaidPhoneAddress"), dataT)
			if err1 != nil {
				ReturnErr101(c, "sorry  system is  wrong4")
				return
			}

			var jsonData PlatformOrderReturnData
			err = json.Unmarshal([]byte(rd), &jsonData)
			if err != nil {
				ReturnErr101(c, "sorry  system is  wrong5")
				return
			}
			if jsonData.Code != 200 {
				ReturnErr101(c, jsonData.Msg)
				return
			}

			//解密
			//base64=> []byte
			decodeString, err1 := base64.StdEncoding.DecodeString(jsonData.Result.(string))
			if err1 != nil {
				zap.L().Debug("recharge err2:" + err1.Error())
				ReturnErr101(c, "sorry  system is  wrong6")
				return
			}

			JumpUrl, err := tools.RsaDecryptForEveryOne(decodeString)
			if err != nil {
				zap.L().Debug("recharge err3:" + err1.Error())
				ReturnErr101(c, "sorry  system is  wrong7")
				return
			}
			type UrlAddressR struct {
				UrlAddress  string
			}
			var  uu UrlAddressR
			err = json.Unmarshal(JumpUrl, &uu)
			if err != nil {
				zap.L().Debug("recharge err4:" + err1.Error())
				ReturnErr101(c, "sorry  system is  wrong8")
				return
			}
			ReturnSuccessData(c, uu.UrlAddress, "ok")
			return

		}

	}
	if action == "ZCS" {
		////类型 1充值  2提现   4购买业务 5佣金奖励(邀请奖励) 6充值到余额宝  7任务提成(团队)
		kinds, _ := strconv.Atoi(c.Query("kinds"))
		record := make([]model.Record, 0)
		//kinds  1  收入  2支出  3充值
		if kinds == 1 {
			mysql.DB.Where("kinds=?  or  kinds=? or  kinds=? ", 1, 5, 7).Where("worker_id=?",whoMap["ID"]).Find(&record)
		}
		if kinds == 2 {
			mysql.DB.Where("kinds=?  or  kinds=?  or kinds=?", 2, 4, 6).Where("worker_id=?",whoMap["ID"]).Find(&record)
		}
		if kinds == 3 {
			mysql.DB.Where("kinds=?  ", 1).Where("worker_id=?",whoMap["ID"]).Find(&record)

		}
		ReturnSuccessData(c, record, "success")
		return
	}



}
