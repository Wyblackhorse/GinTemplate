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
		mysql.DB.Where("kinds=?", kinds).Find(&record)
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
		ReturnSuccess(c, "Withdrawal successful, waiting for administrator review")
		return
	}

	//充值
	if action == "recharge" {

	}

}
