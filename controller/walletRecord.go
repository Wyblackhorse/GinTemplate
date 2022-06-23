package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"net/http"
	"strconv"
)

func GetWalletRecord(c *gin.Context) {
	action := c.Query("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		var total int = 0
		Db := mysql.DB
		fish := make([]model.WalletRecord, 0)
		Db.Model(&model.WalletRecord{}).Count(&total)
		Db = Db.Model(&fish).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		if err := Db.Find(&fish).Error; err != nil {
			tools.JsonWrite(c, -101, nil, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": fish,
		})
		return
	}
}
