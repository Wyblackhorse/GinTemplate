package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"net/http"
	"strconv"
)

func GetBettingRecord(c *gin.Context) {
	action := c.Query("action")

	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		var total int = 0
		Db := mysql.DB

		//按照结算时间

		if closeStart, isExist := c.GetQuery("close_start"); isExist == true {
			if closeEnd, isExist := c.GetQuery("close_end"); isExist == true {
				fmt.Println(closeStart, closeEnd)
			}

		}

		//按照正式账号  正常玩家  私自下注玩家
		if status, isExist := c.GetQuery("status"); isExist == true {
			if status == "保本" {
				Db = Db.Where("break_even=?", "正常")
			} else {

				Db = Db.Where("break_even!=?", "正常")
			}
		}

		fish := make([]model.BettingRecord, 0)
		Db.Table("betting_records").Count(&total)
		Db = Db.Model(&model.BettingRecord{}).Offset((page - 1) * limit).Limit(limit)
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
