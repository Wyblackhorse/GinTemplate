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

//获取订单
func GetRecords(c *gin.Context) {
	action := c.Query("action")
	//获取订单
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		kinds := c.Query("kinds")
		role := make([]model.Record, 0)
		Db := mysql.DB
		var total int

		Db = Db.Where("kinds=?", kinds)
		Db = Db.Model(&model.Record{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Table("records").Count(&total)
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
	//提现审核
	if action == "withdrawalAudit" {
		recordId, _ := strconv.Atoi(c.Query("record_id"))
		status, _ := strconv.Atoi(c.Query("status"))
		record := model.Record{ID: uint(recordId)}
		result, _ := record.IsExistRecord(mysql.DB)
		if result == false {
			ReturnErr101(c, "订单不存在")
			return
		}
		if record.Status != 2 || record.Kinds != 2 {
			ReturnErr101(c, "非法修改!")
			return
		}
		//审核 通过
		resultBool, _ := record.WithdrawDeposit(mysql.DB, status)
		if resultBool == false {
			ReturnErr101(c, "审核失败")
			return
		}
		ReturnSuccess(c, "审核成功")
		return

	}

}
