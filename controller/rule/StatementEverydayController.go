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
	"net/http"
	"strconv"
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
		Db = Db.Model(&model.DailyStatistics{}).Offset((page - 1) * limit).Limit(limit).Order("updated desc")
		Db.Table("daily_statistics").Count(&total)
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

//全局统计
func GetGlobalStatistics(c *gin.Context) {

}
