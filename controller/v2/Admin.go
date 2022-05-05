package V2

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
)

func Login(c *gin.Context) {
	admin := model.Admin{}
	err := mysql.DB.Where("username=?", c.Query("username")).Where("password=?", c.Query("password")).First(&admin).Error
	if err != nil {
		tools.ReturnError101(c, "fail")
		return
	}
	tools.ReturnError200Data(c, admin, "success")
	return

}
