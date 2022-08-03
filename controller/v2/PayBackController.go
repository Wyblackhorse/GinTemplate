package V2

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetPayInformationBack 接受订单回调
func GetPayInformationBack(c *gin.Context) {
	var jsonDataTwo ReturnBase64
	err := c.BindJSON(&jsonDataTwo)
	if err != nil {
		tools.ReturnError101(c, "err:"+err.Error())
		return
	}

	var apiKey = viper.GetString("eth.ApiKey")
	if tools.ApiSign([]byte(apiKey), []byte(jsonDataTwo.Data), []byte(apiKey)) != jsonDataTwo.Sign {
		tools.ReturnError101(c, "非法请求")
		return
	}
	sDec, err1 := base64.StdEncoding.DecodeString(jsonDataTwo.Data)
	if err1 != nil {
		tools.ReturnError101(c, "非法请求")
		return
	}

	//fmt.Println(string(sDec))
	var jsonData GetPayInformationBackData
	err = json.Unmarshal(sDec, &jsonData)
	if err != nil {
		tools.ReturnError101(c, "非法请求")
		return
	}

	if jsonData.Type == "balance" {
		var jsonDataTwo BalanceType
		err = json.Unmarshal(sDec, &jsonDataTwo)
		if err != nil {
			tools.ReturnError101(c, "非法请求")
			return
		}
		zap.L().Debug("余额变动,用户:" + jsonDataTwo.Data.User)
		re := model.ReceiveAddress{Username: jsonDataTwo.Data.User}
		re.UpdateReceiveAddressLastInformationTo0(mysql.DB)
		tools.ReturnError200(c, "余额变动成功")
		return
	}
	p := model.PayOrder{}
	p.TxHash = jsonData.Data.TxHash
	if p.IfIsExitsThisData(mysql.DB) {
		tools.ReturnError200(c, "不要重复添加")
		return
	}
	//添加
	p.Token = strings.ToUpper(jsonData.Data.Token)

	if p.Token == "USDT" {
		acc := strconv.Itoa(jsonData.Data.Amount)
		p.Amount, _ = tools.ToDecimal(acc, 6).Float64()
	} else if p.Token == "TRX" {

	}
	p.FromAddress = jsonData.Data.From
	p.ToAddress = jsonData.Data.To
	p.UserID = jsonData.Data.UserID
	p.Created = time.Now().Unix()
	p.BlockNumber = jsonData.Data.BlockNumber
	p.Timestamp = jsonData.Data.Timestamp / 1000
	p.Date = time.Now().Format("2006-01-02")
	err = mysql.DB.Save(&p).Error
	if err != nil {
		tools.ReturnError101(c, "插入失败:"+err.Error())
		return
	}


	//寻找这个账号最早的充值订单
	p1 := model.PrepaidPhoneOrders{Username: p.UserID, Successfully: p.Timestamp, AccountPractical: p.Amount, RechargeType: strings.ToUpper(p.Token), RechargeAddress: p.FromAddress, CollectionAddress: p.ToAddress}
	p1.UpdateMaxCreatedOfStatusToTwo(mysql.DB, viper.GetInt64("eth.OrderEffectivityTime"))

	//更新钱包地址
	newMoney, _ := tools.ToDecimal(jsonData.Data.Balance, 6).Float64()
	R := model.ReceiveAddress{LastGetAccount: p.Amount, Username: p.UserID, Updated: p.Timestamp, Money: newMoney}
	R.UpdateReceiveAddressLastInformation(mysql.DB)
	da := model.DailyStatistics{RechargeAccount: p.Amount}
	da.UpdateDailyStatistics(mysql.DB)

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

		// 用户名
		if content, isExist := c.GetQuery("UserID"); isExist == true {
			Db = Db.Where("user_id=?", content)
		}
		//	From        string //转账地址
		if content, isExist := c.GetQuery("From"); isExist == true {
			Db = Db.Where("from_Address=?", content)
		}
		//ToAddress
		if content, isExist := c.GetQuery("ToAddress"); isExist == true {
			Db = Db.Where("to_address=?", content)
		}
		//日期条件
		if start, isExist := c.GetQuery("start_time"); isExist == true {
			if end, isExist := c.GetQuery("end_time"); isExist == true {
				Db = Db.Where("timestamp >= ?", start).Where("timestamp<=?", end)
			}
		}
		Db.Table("pay_orders").Count(&total)
		Db = Db.Model(&model.PayOrder{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
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
