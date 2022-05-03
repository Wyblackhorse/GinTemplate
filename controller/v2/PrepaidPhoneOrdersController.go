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

// CreatePrepaidPhoneOrders 生成订单(前端传过来了)
func CreatePrepaidPhoneOrders(c *gin.Context) {
	var jsonData CreatePrepaidPhoneOrdersData
	err := c.BindJSON(&jsonData)
	if err != nil {
		tools.ReturnError101(c, "err:"+err.Error())
		return
	}
	//判断是否存在这个这个用户
	re := model.ReceiveAddress{Username: jsonData.Username, Address: jsonData.CollectionAddress}
	if !re.ReceiveAddressIsExits(mysql.DB) {
		//不存在这个用户 首先要创建这个用户
		re.CreateUsername(mysql.DB, viper.GetString("eth.ThreeUrl"))
	}

	//生成充值订单
	p := model.PrepaidPhoneOrders{ThreeOrder: jsonData.ThreeOrder, CollectionAddress: jsonData.CollectionAddress, AccountOrders: jsonData.AccountOrders, Username: jsonData.Username, RechargeType: jsonData.RechargeType}
	_, err = p.CreatePrepaidPhoneOrders(mysql.DB)
	if err != nil {
		tools.ReturnError101(c, err.Error())
		return
	}

	//充值订单创建成功
	tools.ReturnError200(c, "订单充值成功")
	return
}

func GetPrepaidPhoneOrders(c *gin.Context) {
	action := c.Query("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.PrepaidPhoneOrders, 0)
		Db := mysql.DB
		var total int
		Db = Db.Model(&model.PrepaidPhoneOrders{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Table("prepaid_phone_orders").Count(&total)
		err := Db.Find(&role).Error
		if err != nil {
			tools.ReturnError101(c, "ERR:"+err.Error())
			return
		}

		for k, v := range role {
			address := model.ReceiveAddress{}
			err := mysql.DB.Where("username=?", v.Username).First(&address).Error
			if err == nil {
				role[k].RechargeAddress = address.Address
			}

		}

		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"count": total,
			"data":  role,
		})
		return
	}
}

func Getaddr(c *gin.Context) {
	//{"error":"0","message":"","result":"4564554545454545"}
	c.JSON(http.StatusOK, gin.H{
		"error":   "0",
		"message": "",
		"result":  "TW2HWaLWy9pwiRN4yLju6YKW3aQ6Fw8888",
	})
	return
}
