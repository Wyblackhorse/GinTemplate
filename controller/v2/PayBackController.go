package v2

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"net/http"
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
		tools.ReturnError200(c, "不要重复添加")
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
	p.Timestamp = jsonData.Timestamp / 1000
	err = mysql.DB.Save(&p).Error
	if err != nil {
		tools.ReturnError101(c, "插入失败:"+err.Error())
		return
	}

	//寻找这个账号最早的充值订单
	p1 := model.PrepaidPhoneOrders{Username: p.UserID, Successfully: p.Timestamp, AccountPractical: p.Amount}
	result := p1.UpdateMaxCreatedOfStatusToTwo(mysql.DB)
	if result {
		//更新钱包地址
		R := model.ReceiveAddress{LastGetAccount: p.Amount, Username: p.UserID}
		R.UpdateReceiveAddressLastInformation(mysql.DB)
	}
	tools.ReturnError200(c, "插入成功")
	return

}

// GetPayInformation 获取数据
func GetPayInformation(c *gin.Context) {
	action := c.Query("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.PayOrder, 0)
		Db := mysql.DB
		var total int
		Db = Db.Model(&model.PayOrder{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Table("pay_orders").Count(&total)
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

}
