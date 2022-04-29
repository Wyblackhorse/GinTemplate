/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package rule

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"strconv"
)

//系统设置
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
		fmt.Println(updateMap)

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
