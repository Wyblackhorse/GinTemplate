package v2

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"strconv"
	"time"
)

// GetPayInformationBack 接受订单回调
func GetPayInformationBack(c *gin.Context) {
	var jsonData GetPayInformationBackData
	err := c.BindJSON(&jsonData)
	if err != nil {
		tools.ReturnError101(c, "err:"+err.Error())
		return
	}
	p := model.PayOrder{}
	p.TxHash = jsonData.TxHash
	if p.IfIsExitsThisData(mysql.DB) {
		tools.ReturnError101(c, "不要重复添加")
		return
	}
	//添加
	p.Token = jsonData.Token

	if p.Token == "usdt" {
		acc := strconv.Itoa(jsonData.Amount)
		p.Amount = int(tools.ToDecimal(acc, 6).IntPart())
	}
	p.From = jsonData.From
	p.ToAddress = jsonData.To
	p.UserID = jsonData.UserID
	p.Created = time.Now().Unix()
	p.BlockNumber = jsonData.BlockNumber
	p.Timestamp = jsonData.Timestamp
	err = mysql.DB.Save(&p).Error
	if err != nil {
		tools.ReturnError101(c, "插入失败:"+err.Error())
		return
	}
	tools.ReturnError200(c, "插入成功")
	return

}
