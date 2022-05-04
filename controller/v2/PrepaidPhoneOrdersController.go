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
	re := model.ReceiveAddress{Username: jsonData.Username}
	if !re.ReceiveAddressIsExits(mysql.DB) {
		//不存在这个用户 首先要创建这个用户
		re.CreateUsername(mysql.DB, viper.GetString("eth.ThreeUrl"))
	}

	//生成充值订单
	p := model.PrepaidPhoneOrders{PlatformOrder: jsonData.PlatformOrder, RechargeAddress: jsonData.RechargeAddress, AccountOrders: jsonData.AccountOrders, Username: jsonData.Username, RechargeType: jsonData.RechargeType}
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

		// 用户名
		if content, isExist := c.GetQuery("Username"); isExist == true {
			Db = Db.Where("username=?", content)
		}

		//平台订单号
		if content, isExist := c.GetQuery("PlatformOrder"); isExist == true {
			Db = Db.Where("platform_order=?", content)
		}

		//三方平台订单号
		if content, isExist := c.GetQuery("ThreeOrder"); isExist == true {
			Db = Db.Where("three_order=?", content)
		}

		//充值地址
		if content, isExist := c.GetQuery("RechargeAddress"); isExist == true {
			Db = Db.Where("recharge_address=?", content)
		}

		//订单状态
		if content, isExist := c.GetQuery("Status"); isExist == true {
			Db = Db.Where("status=?", content)
		}

		//是否回调
		if content, isExist := c.GetQuery("ThreeBack"); isExist == true {
			Db = Db.Where("three_back=?", content)
		}

		//日期条件
		if start, isExist := c.GetQuery("start_time"); isExist == true {
			if end, isExist := c.GetQuery("end_time"); isExist == true {
				Db = Db.Where("successfully >= ?", start).Where("successfully<=?", end)
			}
		}

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
				role[k].CollectionAddress = address.Address
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
