/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package worker

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
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
		mysql.DB.Where("kinds=? and worker_id=?", kinds,whoMap["ID"]).Order("created desc").Find(&record)
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

		//d := model.DailyStatistics{TodayRechargeMoney: money,TodayRechargePeople: 1,WorkerId: int(whoMap["ID"].(uint))}
		//d.SetEverydayData(mysql.DB)
	}
	if action == "ZCS" {
		////类型 1充值  2提现   4购买业务 5佣金奖励(邀请奖励) 6充值到余额宝  7任务提成(团队)
		kinds, _ := strconv.Atoi(c.Query("kinds"))
		record := make([]model.Record, 0)
		//kinds  1  收入  2支出  3充值
		if kinds == 1 {
			mysql.DB.Where("kinds=?  or  kinds=? or  kinds=?", 1, 5, 7).Find(&record)
		}
		if kinds == 2 {
			mysql.DB.Where("kinds=?  or  kinds=?  or kinds=?", 2, 4, 6).Find(&record)
		}
		if kinds == 3 {
			mysql.DB.Where("kinds=?  ", 1).Find(&record)

		}
		ReturnSuccessData(c, record, "success")
		return
	}
}


