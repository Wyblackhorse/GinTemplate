/**
 * @Author $
 * @Description //TODO $ 报表数据
 * @Date $ $
 * @Param $
 * @return $
 **/
package rule

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"net/http"
	"strconv"
	"time"
)

//每日报表
func GetStatementEveryday(c *gin.Context) {
	action := c.Query("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.DailyStatistics, 0)
		Db := mysql.DB
		var total int
		Db.Table("daily_statistics").Count(&total)

		Db = Db.Model(&model.DailyStatistics{}).Offset((page - 1) * limit).Limit(limit).Order("updated desc")
		err := Db.Find(&role).Error
		if err != nil {
			ReturnErr101(c, "ERR:"+err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": role,
		})
		return
	}
}

//团队列表
func GetTeamStatistics(c *gin.Context) {
	action := c.Query("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.Worker, 0)
		Db := mysql.DB
		var total int
		Db.Table("workers").Count(&total)
		Db = Db.Model(&model.Worker{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		err := Db.Find(&role).Error
		if err != nil {
			ReturnErr101(c, "ERR:"+err.Error())
			return
		}

		type ReturnData struct {
			Phone          string  //账号
			PeopleNum      int     `json:"people_num"`      //团队人数
			Recharge       float64 `json:"recharge"`        //充值
			Withdraw       float64 `json:"withdraw"`        //提现
			TaskNum        int     `json:"task_num"`        //任务数量
			PromotionAward float64 `json:"promotion_award"` //推广奖励
			RevocationTask int     `json:"revocation_task"` //撤销任务
			TaskPercentage float64 `json:"task_percentage"` //任务提成
			ids            []uint  `json:"ids"`
		}

		re := make([]ReturnData, 0)
		for _, v := range role {
			var u ReturnData
			u.Phone = v.Phone
			//团队人数
			mysql.DB.Model(&model.Worker{}).Where("superior_id=? or  next_superior_id=? or next_next_superior_id=?", v.ID, v.ID, v.ID).Count(&u.PeopleNum)

			o := make([]model.Worker, 0)
			mysql.DB.Model(&model.Worker{}).Where("superior_id=? or  next_superior_id=? or next_next_superior_id=?", v.ID, v.ID, v.ID).Find(&o)

			var ids []int
			for _, i2 := range o {
				ids = append(ids, int(i2.ID))
			}
			//充值金额
			mysql.DB.Raw("SELECT SUM(money) as recharge from records where status=1 and kinds=1 and worker_id in (?)", ids).Scan(&u)
			//提现
			mysql.DB.Raw("SELECT SUM(money) as withdraw from records where status=1 and kinds=2 and worker_id in (?)", ids).Scan(&u)
			////推广奖励
			mysql.DB.Raw("SELECT SUM(money) as promotion_award from records where status=1 and kinds=5 and worker_id in (?)", ids).Scan(&u)
			////任务数量
			mysql.DB.Raw("SELECT count(*)  as task_num from task_orders where status=3 and worker_id in (?)", ids).Scan(&u)
			////撤销的任务(失败的)
			mysql.DB.Raw("SELECT count(*)  as revocation_task from task_orders where status=4 and  worker_id in (?)", ids).Scan(&u)
			////任务提成
			mysql.DB.Raw("SELECT SUM(money) as task_percentage from records where status=1 and kinds=5 and worker_id in (?)", ids).Scan(&u)
			re = append(re, u)

		}
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": re,
		})
		return

	}

}

//全局统计
func GetGlobalStatistics(c *gin.Context) {
	type GetGlobalStatistics struct {
		Today     model.DailyStatistics
		Yesterday model.DailyStatistics
		ThisWeek  model.DailyStatistics
		ThisMonth model.DailyStatistics
	}
	//今日统计
	var dd GetGlobalStatistics
	mysql.DB.Where("date=?", time.Now().Format("2006-01-02")).First(&dd.Today)
	//昨日
	mysql.DB.Where("date=?", time.Now().AddDate(0, 0, -1).Format("2006-01-02")).First(&dd.Yesterday)
	//这周
	mysql.DB.Where("week=?", tools.ReturnTheWeek()).First(&dd.ThisWeek)
	//这个月
	mysql.DB.Where("month=?", tools.ReturnTheMonth()).First(&dd.ThisMonth)
	ReturnSuccessData(c, dd, "success")
	return

}
