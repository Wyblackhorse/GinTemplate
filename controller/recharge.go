package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"net/http"
	"strconv"
	"time"
)

func GetRecharge(c *gin.Context) {
	action := c.Query("action")

	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		var total int = 0
		Db := mysql.DB
		fish := make([]model.Recharge, 0)
		Db.Table("recharges").Count(&total)
		Db = Db.Model(&fish).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		if err := Db.Find(&fish).Error; err != nil {
			tools.JsonWrite(c, -101, nil, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": fish,
		})
		return
	}

	if action == "yesterday" {
		//冲提报表  查看昨天充提报表全部复制下来
		type data struct {
			Username          string  `json:"username"`             //玩家
			TheGeneralAgentOf string  `json:"the_general_agent_of"` //总代理用户名
			Recharges         float64 `json:"recharges"`            //充值金额
			WithdrawalAmount  float64 `json:"withdrawal_amount"`    //提现金额

		}
		fmt.Println(time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
		var aa []data
		mysql.DB.Raw("SELECT app_users.username  as username,app_users.the_general_agent_of as  the_general_agent_of, recharges.recharge  as  recharges,withdraws.withdrawal_amount as  withdrawal_amount FROM recharges LEFT JOIN withdraws ON withdraws.the_user_id=recharges.user_id LEFT JOIN app_users ON app_users.user_number=recharges.user_id where recharges.date=? ", time.Now().AddDate(0, 0, -1).Format("2006-01-02")).Scan(&aa)
		tools.JsonWrite(c, 200, aa, "ok")
		return
	}

	//统计用户充值和提现
	if action == "TotalWAndR" {

		type data struct {
			Username          string  `json:"username"`             //玩家
			TheGeneralAgentOf string  `json:"the_general_agent_of"` //总代理用户名
			Recharges         float64 `json:"recharges"`            //充值金额
			WithdrawalAmount  float64 `json:"withdrawal_amount"`    //提现金额

		}
		var aa []data
		mysql.DB.Raw("SELECT   app_users.username  as username ,app_users.the_general_agent_of as  the_general_agent_of, sum(recharges.recharge) as recharges ,SUM(withdraws.withdrawal_amount) as withdrawal_amount FROM recharges LEFT JOIN withdraws ON withdraws.the_user_id=recharges.user_id   LEFT JOIN app_users ON app_users.user_number=recharges.user_id  WHERE recharges.`status`='成功'  OR withdraws.`status`='通过' GROUP BY  app_users.username ").Scan(&aa)
		tools.JsonWrite(c, 200, aa, "ok")
		return

	}

}
