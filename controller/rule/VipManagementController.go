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

//获取会员等级  修改 删除
func GetVipLevel(c *gin.Context) {
	action := c.Query("action")
	//获取vip
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.Vip, 0)
		Db := mysql.DB
		var total int
		Db = Db.Model(&model.Vip{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Table("vips").Count(&total)
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

	//修改
	if action == "UPDATE" {

	}

	if action == "ADD" {

	}

	if action == "DEL" {

	}
	ReturnErr101(c, "no action")
	return

}

//获取会员成员
func GetVipWorkers(c *gin.Context) {

	action := c.Query("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.Worker, 0)
		Db := mysql.DB
		var total int
		Db = Db.Model(&model.Worker{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Table("workers").Count(&total)
		err := Db.Find(&role).Error
		if err != nil {
			ReturnErr101(c, "ERR:"+err.Error())
			return
		}
		for k, v := range role {
			//判断会员等级
			vip := model.Vip{ID: uint(v.VipId)}
			role[k].VipName = vip.GetLevelName(mysql.DB)
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": role,
		})
		return
	}

	//更新用户
	if action == "UPDATE" {
		workerId, _ := strconv.Atoi(c.Query("worker_id"))
		worker := model.Worker{ID: uint(workerId)}
		if !worker.IsExist(mysql.DB) {
			ReturnErr101(c, "这个用户不存在")
			return
		}
		update := map[string]interface{}{}
		if updateData, isExist := c.GetQuery("status"); isExist == true {
			update["Status"], _ = strconv.Atoi(updateData)
		}
		err := mysql.DB.Model(model.Worker{}).Where("id=?", workerId).Update(update).Error
		if err != nil {
			ReturnErr101(c, "修改失败")
			return
		}
		ReturnErr101(c, "修改成功")
		return

	}

	ReturnErr101(c, "no action")
	return
}
