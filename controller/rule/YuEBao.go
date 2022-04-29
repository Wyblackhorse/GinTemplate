/**
 * @Author $
 * @Description //TODO $
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

//获取余额包 活动列表
func GetYuEBaoList(c *gin.Context) {
	action, _ := c.Get("action")
	//获取活动数据
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.YuEBao, 0)
		Db := mysql.DB
		var total int
		Db = Db.Model(&model.YuEBao{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Table("yu_e_baos").Count(&total)
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
	//添加数据
	if action == "ADD" {
		name := c.Query("name")
		minMoney, _ := strconv.ParseFloat(c.Query("min_money"), 64)
		interestRate, _ := strconv.ParseFloat(c.Query("interest_rate"), 64)
		Days, _ := strconv.Atoi(c.Query("days"))
		bao := model.YuEBao{Name: name, MinMoney: minMoney, InterestRate: interestRate, Days: Days}
		result, _ := bao.AddYuEBao(mysql.DB)
		if result == false {
			ReturnErr101(c, "添加失败")
			return
		}
		ReturnSuccess(c, "添加成功")
		return
	}

}
