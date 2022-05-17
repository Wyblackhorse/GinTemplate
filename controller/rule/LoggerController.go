/**
 * @Author $
 * @Description //TODO $日志管理
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

//管理操作日志
func GetAdminLog(c *gin.Context) {
	action := c.Query("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		//WriteId int    //操作者  如果为0  就是系统日志
		//Kinds   int    //1玩家 2管理员
		Apply := make([]model.Loggers, 0)
		Db := mysql.DB
		if kinds, isExist := c.GetQuery("kinds"); isExist == true {
			Db = Db.Where("kinds=?", kinds)
		}
		if kinds, isExist := c.GetQuery("WriteId"); isExist == true {
			Db = Db.Where("write_id=?", kinds)
		}
		if kinds, isExist := c.GetQuery("kinds"); isExist == true {
			Db = Db.Where("kinds=?", kinds)
		}
		var total int
		Db = Db.Model(&model.Loggers{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Table("loggers").Count(&total)
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
}

//获取管理者
func GetAdmin(c *gin.Context) {
	action := c.Query("action")
	if action == "GET" {
		type AdminModel struct {
			ID       uint
			Username string
		}
		//mysql.DB.Table("admin_models").Find(&results)
		admins := make([]AdminModel, 0)
		mysql.DB.Find(&admins)
		ReturnSuccessData(c, admins, "获取成功")
		return
	}
}

//系统日志

func GetSystemLog(c *gin.Context) {
	action := c.Query("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		//WriteId int    //操作者  如果为0  就是系统日志
		//Kinds   int    //1玩家 2管理员
		Apply := make([]model.Loggers, 0)
		Db := mysql.DB
		Db = Db.Where("write_id=?", 0)
		var total int
		Db = Db.Model(&model.Loggers{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Table("loggers").Count(&total)
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
}
