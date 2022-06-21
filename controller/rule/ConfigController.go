/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package rule

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"net/http"
	"strconv"
	"strings"
)

//系统设置   基本设置
func SetConfig(c *gin.Context) {
	action := c.Query("action")
	if action == "GET" {
		config := model.Config{}
		mysql.DB.Where("id=?", 1).First(&config)
		ReturnSuccessData(c, config, "获取成功")
		return
	}
	if action == "UPDATE" {
		updateMap := map[string]interface{}{}
		//是否需要验证码
		if add, isExist := c.GetQuery("need_invitation_code"); isExist == true {
			updateMap["NeedInvitationCode"], _ = strconv.Atoi(add)
		}
		//默认信用分
		if add, isExist := c.GetQuery("credit_score"); isExist == true {
			updateMap["CreditScore"], _ = strconv.Atoi(add)
		}
		//邀请奖励
		if add, isExist := c.GetQuery("invite_rewards"); isExist == true {
			updateMap["InviteRewards"], _ = strconv.ParseFloat(add, 64)
		}
		//上级返点
		if add, isExist := c.GetQuery("superior_back_per"); isExist == true {
			updateMap["SuperiorBackPer"], _ = strconv.ParseFloat(add, 64)
		}
		//次上级返点
		if add, isExist := c.GetQuery("next_superior_back_per"); isExist == true {
			updateMap["NextSuperiorBackPer"], _ = strconv.ParseFloat(add, 64)
		}
		//次次上级返点
		if add, isExist := c.GetQuery("next_next_superior_back_per"); isExist == true {
			updateMap["NextNextSuperiorBackPer"], _ = strconv.ParseFloat(add, 64)
		}
		//最低提现金额
		if add, isExist := c.GetQuery("low_withdrawal"); isExist == true {
			updateMap["LowWithdrawal"], _ = strconv.ParseFloat(add, 64)
		}
		//提现手续费
		if add, isExist := c.GetQuery("low_withdrawal"); isExist == true {
			updateMap["WithdrawalCharge"], _ = strconv.ParseFloat(add, 64)
		}
		//开通云管家最低等级
		if add, isExist := c.GetQuery("open_cloud_housekeeper_level_id"); isExist == true {
			updateMap["OpenCloudHousekeeperLevelId"], _ = strconv.Atoi(add)
		}
		//前台默认  语言
		if add, isExist := c.GetQuery("ForegroundLanguage"); isExist == true {
			updateMap["ForegroundLanguage"] = add
		}
		///完成当日首次任务加用分
		if add, isExist := c.GetQuery("DoneFirstTask"); isExist == true {
			updateMap["DoneFirstTask"], _ = strconv.Atoi(add)
		}
		//没有在规定时间完成任务扣除信用分
		if add, isExist := c.GetQuery("OverTimeForTask"); isExist == true {
			updateMap["OverTimeForTask"], _ = strconv.Atoi(add)
		}
		//没有按照要求完成任务扣除信用分
		if add, isExist := c.GetQuery("NoRequireToTask"); isExist == true {
			updateMap["NoRequireToTask"], _ = strconv.Atoi(add)
		}
		//推荐奖励次数     1不重复发放  2重复发放
		if add, isExist := c.GetQuery("ReferralBonusesTimes"); isExist == true {
			updateMap["ReferralBonusesTimes"], _ = strconv.Atoi(add)
		}

		//最高提现金额
		if add, isExist := c.GetQuery("HighWithdrawal"); isExist == true {
			updateMap["HighWithdrawal"], _ = strconv.ParseFloat(add, 64)
		}

		//客服地址
		if add, isExist := c.GetQuery("LinkOfTheService"); isExist == true {
			updateMap["LinkOfTheService"] = add
		}
		//飞机地址
		if add, isExist := c.GetQuery("TelegramAddress"); isExist == true {
			updateMap["TelegramAddress"] = add
		}
		err := mysql.DB.Model(&model.Config{}).Where("id=?", 1).Update(updateMap).Error
		if err != nil {
			ReturnErr101(c, "更新失败")
			return
		}
		ReturnSuccess(c, "更新成功")
		return
	}
	ReturnErr101(c, "no fond action")
	return
}

//幻灯片设置
func LanternSlide(c *gin.Context) {
	action := c.Query("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.LanternSlide, 0)
		Db := mysql.DB
		var total int

		//语言
		if la, isExist := c.GetQuery("language"); isExist == true {
			Db = Db.Where("language=?", la)
		}

		//状态
		if la, isExist := c.GetQuery("status"); isExist == true {
			status, _ := strconv.Atoi(la)
			Db = Db.Where("status=?", status)
		}

		Db = Db.Model(&model.DailyStatistics{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Table("lantern_slides").Count(&total)
		err := Db.Find(&role).Error
		if err != nil {
			ReturnErr101(c, "ERR:"+err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": role,
		})
		return
	}
	// 添加
	if action == "ADD" {
		file, err := c.FormFile("file")
		if err != nil {
			ReturnErr101(c, err.Error())
			return
		}
		fileArray := strings.Split(file.Filename, ".")
		if fileArray[1] != "png" && fileArray[1] != "jpg" {
			ReturnErr101(c, "上传文件格式不对,只接受 *.png  *.jpg")
			return
		}
		ua := "./static/slideshow/" + file.Filename
		err = c.SaveUploadedFile(file, ua)
		if err != nil {
			ReturnErr101(c, err.Error())
			return
		}
		err2 := mysql.DB.Where("url_address=? and language=?", ua, c.PostForm("language")).First(&model.LanternSlide{}).Error
		if err2 == nil {
			ReturnErr101(c, "不要重复添加")
			return
		}
		slide := model.LanternSlide{Language: c.PostForm("language"), UrlAddress: ua}
		if slide.Add(mysql.DB) == false {
			ReturnErr101(c, "添加失败")
			return
		}
		ReturnSuccess(c, "上传成功")
		return
	}
}
