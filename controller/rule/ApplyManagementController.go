/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $   应用管理
 **/
package rule

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"net/http"
	"strconv"
)

//获取应用列表
func GetApplyList(c *gin.Context) {
	action := c.Query("action")
	//获取
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		Apply := make([]model.Apply, 0)
		Db := mysql.DB
		var total int
		Db = Db.Model(&model.Apply{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Table("roles").Count(&total)
		err := Db.Find(&Apply).Error
		if err != nil {
			ReturnErr101(c, "ERR:"+err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": Apply,
		})
		return
	}
	//更新
	if action == "UPDATE" {

	}
	//删除
	if action == "DEL" {

	}
	//添加
	if action == "ADD" {

	}

}
