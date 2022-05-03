package v2

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"net/http"
	"strconv"
)

// GetReceiveAddress 获取地址管理
func GetReceiveAddress(c *gin.Context) {
	action := c.Query("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.ReceiveAddress, 0)
		Db := mysql.DB
		var total int
		Db = Db.Model(&model.ReceiveAddress{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Table("receive_addresses").Count(&total)
		err := Db.Find(&role).Error
		if err != nil {
			tools.ReturnError101(c, "ERR:"+err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"count": total,
			"data":  role,
		})
		return
	}

	if action == "ADD" {
		result := model.CreateNewReceiveAddress(mysql.DB, viper.GetString("eth.ThreeUrl"))
		if !result {
			tools.ReturnError101(c, "添加失败")
			return
		}
		tools.ReturnError200(c, "添加成功")
		return
	}

}
