/**
 * @Author $
 * @Description //TODO 理财产品
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

//购买理财产品
func ShoppingMoneyManagement(c *gin.Context) {

	action := c.Query("action")
	who, _ := c.Get("who")
	whoMap := who.(map[string]interface{})
	//获取理财产品
	if action == "GET" {
		YEB := make([]model.YuEBao, 0)
		err := mysql.DB.Where("status=?", 1).Find(&YEB).Error
		if err != nil {
			ReturnSuccessData(c, []string{}, "success")

			return
		}
		ReturnSuccessData(c, YEB, "success")
		return
	}
	//获取已经自己已经购买的理财产品
	if action == "alreadyBought" {
		moneyManagement := make([]model.MoneyManagement, 0)
		mysql.DB.Where("worker_id=?", whoMap["ID"]).Find(&moneyManagement)
		ReturnSuccessData(c, moneyManagement, "success")
		return
	}
	//购买理财产品
	if action == "shopping" {
		money, err := strconv.ParseFloat(c.Query("money"), 64)
		if err != nil {
			ReturnErr101(c, "fail")
			return
		}

		id, _ := strconv.Atoi(c.Query("id"))

		//判断这个产品是否存在
		man := model.YuEBao{}
		err = mysql.DB.Where("id=?", id).First(&man).Error
		if err != nil {
			ReturnErr101(c, "fail")
			return
		}
		if man.Status == 2 {
			ReturnErr101(c, "closed")
			return
		}
		//判断钱是否满足需求
		if money < man.MinMoney {
			ReturnErr101(c, "don't have enough money")
			return
		}

		//钱够了 进行下面的逻辑

		Worker := model.WorkerBalance{ID: int(whoMap["ID"].(uint)), AddBalance: money, Kinds: 6, YuEBaoId: int(man.ID), Days: man.Days}
		_, _ = Worker.AddBalanceFuc(mysql.DB)

	}

}
