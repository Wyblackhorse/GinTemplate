package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"net/http"
	"strconv"
	"strings"
)

func GetUserApp(c *gin.Context) {
	action := c.Query("action")

	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		var total int = 0
		Db := mysql.DB
		fish := make([]model.AppUser, 0)
		Db.Table("app_users").Count(&total)
		Db = Db.Model(&fish).Offset((page - 1) * limit).Limit(limit).Order("updated desc")
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

	if action == "GetFather" {
		username := c.PostForm("username")
		usernameArray := strings.Split(username, "\n")
		var ReturnData []string
		for _, s := range usernameArray {
			u := model.AppUser{}
			err := mysql.DB.Where("username=?", s).First(&u).Error
			if err == nil {
				ReturnData = append(ReturnData, s+"==>"+u.TheGeneralAgentOf)
			}

		}
		tools.JsonWrite(c, 200, ReturnData, "OK")
		return
	}

	if action == "Details" {
		username := c.Query("username")
		//本级
		a := model.AppUser{}
		err := mysql.DB.Where("username=?", username).First(&a).Error
		if err != nil {
			tools.JsonWrite(c, -101, nil, "用户不存在")
			return
		}
		array := model.DetailAppUserArray{}
		array.Myself = ReturnDetailAppUser(username, a)
		//上级
		a1:=model.AppUser{}
		if a.UpperLayerUserName != "" {
			err = mysql.DB.Where("username=?", a.UpperLayerUserName).First(&a1).Error
			if err != nil {
				tools.JsonWrite(c, 200, array, "OK")
				return
			}
			array.One = ReturnDetailAppUser(a.UpperLayerUserName, a1)
		}


		//上上级
		a2:=model.AppUser{}
		err = mysql.DB.Where("username=?", a1.UpperLayerUserName).First(&a2).Error
		if err != nil {
			tools.JsonWrite(c, 200, array, "OK")
			return
		}

		if a.UpperLayerUserName != "" {
			array.Two = ReturnDetailAppUser(a.UpperLayerUserName, a2)
		}

		//上上级
		a3:=model.AppUser{}
		err = mysql.DB.Where("username=?", a2.UpperLayerUserName).First(&a3).Error
		if err != nil {
			tools.JsonWrite(c, 200, array, "OK")
			return
		}
		if a.UpperLayerUserName != "" {
			array.Three = ReturnDetailAppUser(a.UpperLayerUserName, a3)
		}
		tools.JsonWrite(c, 200, array, "OK")

	}

}

func ReturnDetailAppUser(username string, a model.AppUser) model.DetailAppUser {

	de := model.DetailAppUser{UserNumber: a.UserNumber, UpperLayerUserName: username, RegistrationTime: a.RegistrationTime, LastLoginTime: a.LastLoginTime}
	//余额
	wr := model.WalletRecord{}
	err := mysql.DB.Where("user_name=?", username).First(&wr).Error
	if err == nil {
		de.Money = wr.Amount
	}
	//充值
	rec := model.Recharge{}
	err = mysql.DB.Where("username=? and  status=?", username, "成功").First(&rec).Error
	if err == nil {
		de.Recharge = rec.Recharge
	}
	//提现 the_user_name
	wth := model.Withdraw{}
	err = mysql.DB.Where("the_user_name=? and  status=?", username, "通过").First(&wth).Error
	if err == nil {
		de.Withdraw = wth.WithdrawalAmount
	}
	//直属下级个数
	mysql.DB.Model(&model.AppUser{}).Where("upper_layer_user_name=?", username).Count(&de.DirectlySubordinateNum)
	//团队人数
	de.TeamNums = GetTeamNums(username)
	//上级列表
	de.SuperiorList = GetUpList(a.UpperLayerUserName)
	//下注金额
	mysql.DB.Raw("SELECT sum(bet_amount) as bet_money FROM  betting_records WHERE user_name= ?  and  break_even=?  ", username, "正常").Scan(&de)
	//下注次数
	mysql.DB.Model(&model.BettingRecord{}).Where("user_name=? and break_even=? ", username, "正常").Count(&de.BetNum)
	de.SaveDifference = de.Recharge - de.Withdraw
	return de
}

func GetTeamNums(username string) int {

	app := make([]model.AppUser, 0)
	//第一级
	mysql.DB.Where("upper_layer_user_name=?", username).Find(&app)
	allNums := len(app)
	for _, user := range app {
		add := make([]model.AppUser, 0)
		mysql.DB.Where("upper_layer_user_name=?", user.Username).Find(&add)
		allNums = allNums + len(add)
		for _, i2 := range add {
			adc := make([]model.AppUser, 0)
			mysql.DB.Where("upper_layer_user_name=?", i2.Username).Find(&adc)
			allNums = allNums + len(adc)
		}
	}

	return allNums + 1
}

// GetUpList 获取上级列表
func GetUpList(UpName string) string {
	reD := UpName
	re := model.AppUser{}
	err := mysql.DB.Where("username=?", UpName).First(&re).Error
	if err != nil {
		return reD
	}
	reD = reD + "," + re.UpperLayerUserName


	re2:=model.AppUser{}
	err = mysql.DB.Where("username=?", re.UpperLayerUserName).First(&re2).Error
	if err != nil {
		return reD
	}
	reD = reD + "," + re2.UpperLayerUserName
	re3:=model.AppUser{}
	err = mysql.DB.Where("username=?",  re2.UpperLayerUserName).First(&re3).Error
	if err != nil {
		return reD
	}
	reD = reD + "," + re3.UpperLayerUserName
	return reD
}
