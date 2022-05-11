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

	InvitationCode := register.InvitationCode
	newWork := model.Worker{}
	//需要邀请码
	if config.NeedInvitationCode == 1 && InvitationCode == "" {
		ReturnErr101(c, "The invitation code must not be empty")
		return
	}

	//邀请码为空
	if InvitationCode != "" {
		fmt.Println(InvitationCode)
		//判断邀请码是否存在
		worker := model.Worker{}
		err := mysql.DB.Where("Invitation_code =?", register.InvitationCode).First(&worker).Error
		if err != nil {
			ReturnErr101(c, "The invitation code is not exist")
			return
		}

		if config.InviteRewards > 0 {

			//创建订单
			re := model.Record{WorkerId: int(worker.ID), Status: 1, Money: config.InviteRewards, Kinds: 5}
			reResult, _ := re.AddRecord(mysql.DB)
			if reResult == true {
				//邀请奖励大于0
				w2 := model.WorkerBalance{ID: int(worker.ID), Kinds: 5, AddBalance: config.InviteRewards}
				_, _ = w2.AddBalanceFuc(mysql.DB)
			}

		}

		newWork.SuperiorId = int(worker.ID)                //上级
		newWork.NextSuperiorId = worker.SuperiorId         // 次上级
		newWork.NextNextSuperiorId = worker.NextSuperiorId // 次上上级
	}

	//校验验证码
	verifyResult := store.Verify(register.VerificationId, register.VerificationCode, true)
	if verifyResult != false {
		//验证码校验没有通过
		ReturnErr101(c, "The verify code is err")
		return
	}
	//查看这个账户是否已经注册过了
	err = mysql.DB.Where("phone=?", register.Phone).First(&model.Worker{}).Error
	if err == nil {
		ReturnErr101(c, "The account already exists")
		return
	}

	//

	token := tools.RandStringRunes(60)
	newWork.Password = register.Password
	newWork.Phone = register.Phone
	newWork.Token = token
	newWork.Status = 2
	newWork.InvitationCode = tools.RandStringRunes(8)
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
	return
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
	redis.Rdb.Set("Worker_Login_Token_"+worker.Token, login.Phone, time.Second*3600*24)
	ReturnSuccessData(c, worker, "success")
	return
}

// 获取个人信息
func Information(c *gin.Context) {
	action := c.Query("action")
	who, _ := c.Get("who")
	whoMap := who.(map[string]interface{})
	if action == "GET" {
		ReturnSuccessData(c, whoMap, "success")
		return
	}
	//更新资料
	if action == "UPDATE" {
		update := make(map[string]interface{})
		if username, isExits := c.GetQuery("username"); isExits == true {
			update["Username"] = username
		}
		if username, isExits := c.GetQuery("password"); isExits == true {
			update["Password"] = username
		}
		if username, isExits := c.GetQuery("pay_password"); isExits == true {
			update["PayPassword"] = username
		}
		err := mysql.DB.Where("id=?", whoMap["ID"]).Update(update).Error
		if err != nil {
			ReturnErr101(c, "fail")
			return
		}
		ReturnSuccess(c, "ok")
		return

	}

	//获取收益
	if action == "getBenefit" {
		//returnMap := make(map[string]interface{})
		t := model.TaskOrder{WorkerId: int(whoMap["ID"].(uint))}
		//获取用户的任务输
		vip := model.Vip{}
		fmt.Println(whoMap["VipId"])
		err := mysql.DB.Where("id=?", whoMap["VipId"]).First(&vip).Error
		if err != nil {
			ReturnErr101(c, "no find vip level")
			return
		}
		var res model.ReturnResult
		res = t.GetEarnings(mysql.DB, vip.TaskTimes)
		res.Balance = whoMap["Balance"].(float64)
		ReturnSuccessData(c, res, "OK")
		return
	}

}
