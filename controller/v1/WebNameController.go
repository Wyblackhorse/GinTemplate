package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/dao/redis"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/process"
	"github.com/wangyi/GinTemplate/tools"
	"net/http"
	"strconv"
)

// SetWebName 设置 域名 和获取域名
func SetWebName(c *gin.Context) {
	action := c.Query("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.WebName, 0)
		Db := mysql.DB
		var total int
		Db = Db.Model(&model.WebName{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Table("web_names").Count(&total)
		err := Db.Find(&role).Error
		if err != nil {
			tools.ReturnFail101(c, "ERR:"+err.Error())
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
		url := c.Query("webName")
		u := model.WebName{Url: url}
		_, err := u.AddWebName(mysql.DB)
		if err != nil {
			tools.ReturnFail101(c, err.Error())
			return
		}
		tools.ReturnSuccess200(c, "添加成功")
		return
	}

	//干活
	if action == "DOING" {
		id := c.Query("id")
		web := model.WebName{}
		err := mysql.DB.Where("id=?", id).First(&web).Error
		if err != nil {
			tools.ReturnFail101(c, "不存在的域名")
			return
		}

		//判断是否在任务中
		DONGING, _ := redis.Rdb.HExists("DOING", id).Result()
		if DONGING {
			tools.ReturnFail101(c, "已在任务中,不要重复执行程序")
			return
		}

		go process.SetCheckWebNameIsTrueProcess(web.MatchingUrl, web, redis.Rdb, mysql.DB)
		redis.Rdb.HSet("DOING", id, "doing")
		tools.ReturnSuccess200(c, "执行成功")
		return
	}

}

func CheckWebName(c *gin.Context) {
	action := c.Query("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		status, _ := strconv.Atoi(c.Query("status"))
		role := make([]model.CheckWebName, 0)
		Db := mysql.DB

		Db = Db.Where("status=?", status)

		var total int
		Db = Db.Model(&model.CheckWebName{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Table("check_web_names").Count(&total)
		err := Db.Find(&role).Error

		for i, name := range role {
			W := model.WebName{}
			err := mysql.DB.Where("id=?", name.WebNameId).First(&W).Error
			if err == nil {
				role[i].WebName = W.Url
			}
		}

		if err != nil {
			tools.ReturnFail101(c, "ERR:"+err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"count": total,
			"data":  role,
		})
		return
	}
}
