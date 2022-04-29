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

//获取账单
func GetBill(c *gin.Context)  {


	action := c.Query("action")
	//获取基本数据
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.BillingDetails, 0)
		Db := mysql.DB
		var total int
		Db = Db.Model(&model.BillingDetails{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Table("billing_details").Count(&total)
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
	ReturnErr101(c, "no action")
	return
}
