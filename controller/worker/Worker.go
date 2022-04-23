/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package worker

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/dao/redis"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"time"
)

//注册
func Register(c *gin.Context) {

	var register RegisterRequest
	if err := c.ShouldBind(&register); err != nil {
		tools.JsonWrite(c, -101, []string{}, "err:"+err.Error())
		return
	}

	//判断两次密码是否一样
	if register.Password != register.PasswordTwo {
		ReturnErr101(c, "Inconsistent passwords")
		return
	}

	//查看邀请码是否需要填写?
	config := model.Config{}
	err := mysql.DB.Where("id=?", 1).First(&config).Error
	if err != nil {
		ReturnErr101(c, "fail")
		return
	}

	InvitationCode := register.VerificationCode

	newWork := model.Worker{}
	//需要邀请码
	if config.NeedInvitationCode == 1 {
		if InvitationCode == "" {
			ReturnErr101(c, "The invitation code must not be empty")
			return
		} else {
			//判断邀请码是否存在
			worker := model.Worker{}
			err := mysql.DB.Where("Invitation_code =?", register.InvitationCode).First(&worker).Error
			if err != nil {
				ReturnErr101(c, "The invitation code is not exist")
				return
			}
			newWork.SuperiorId = int(worker.ID)
		}
	}
	//校验验证码
	verifyResult := store.Verify(register.VerificationId, register.VerificationCode, true)
	fmt.Println(verifyResult)
	if verifyResult == false {
		//验证码校验没有通过
		ReturnErr101(c, "The invitation code is err")
		return
	}

	//查看这个账户是否已经注册过了
	err = mysql.DB.Where("phone=?", register.Phone).First(&model.Worker{}).Error
	if err == nil {
		ReturnErr101(c, "The account already exists")
		return
	}
	token := tools.RandStringRunes(40)
	newWork.Password = register.Password
	newWork.Phone = register.Phone
	newWork.Token = token
	newWork.Status = 2
	newWork.CreditScore = config.CreditScore
	newWork.Created = time.Now().Unix()
	err = mysql.DB.Save(&newWork).Error
	if err != nil {
		//注册失败
		ReturnErr101(c, "register is fail")
		return
	}
	//设置每日统计数据
	daily := model.DailyStatistics{Register: 1}
	daily.SetEverydayData(mysql.DB)
	redis.Rdb.HSet("Worker_Token", token, newWork.Phone)
	redis.Rdb.HMSet("Worker_"+newWork.Phone, structs.Map(&newWork))
	ReturnSuccess(c, "注册成功")
}

///登录
func Login(c *gin.Context) {
	//参数校验
	var login LoginRequest
	if err := c.ShouldBind(&login); err != nil {
		tools.JsonWrite(c, -101, []string{}, "err:"+err.Error())
		return
	}
	//判断账户密码是否正确
	worker := model.Worker{}
	err := mysql.DB.Where("phone=?", login.Phone).Where("password=?", login.Password).First(&worker).Error
	if err != nil {
		ReturnErr101(c, "The account or password is incorrect")
		return
	}
	ReturnSuccessData(c, worker, "success")
	return
}
