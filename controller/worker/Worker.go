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

	//查看这个账户是否已经注册过了
	err = mysql.DB.Where("phone=?", register.Phone).First(&model.Worker{}).Error
	if err == nil {
		ReturnErr101(c, "The account already exists")
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
			//创建订单
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
	if verifyResult == false {
		//验证码校验没有通过
		ReturnErr101(c, "The verify code is err")
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
	daily := model.DailyStatistics{TodayRegister: 1}
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

	//账号已被封禁
	if worker.Status == 4 {
		ReturnErr101(c, "The account has since been banned")
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
			if old, isExits := c.GetQuery("old_password"); isExits == true {
				if old != whoMap["Password"] {
					ReturnErr101(c, "The original password entered is incorrect")
					return
				}
				update["Password"] = username
			}
		}
		if username, isExits := c.GetQuery("pay_password"); isExits == true {
			update["PayPassword"] = username
		}
		err := mysql.DB.Model(&model.Worker{}).Where("id=?", whoMap["ID"]).Update(update).Error
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
	//设置 银行卡
	if action == "setBank" {
		name := c.Query("name")
		mail := c.Query("mail")
		phone := c.Query("phone")
		address := c.Query("address")
		if address == "" {
			ReturnErr101(c, "The collection address must not be empty")
			return
		}
		b := model.Bank{Address: address, Name: name, Phone: phone, Mail: mail, WorkerId: int(whoMap["ID"].(uint))}
		_, err := b.Add(mysql.DB)
		if err != nil {
			ReturnErr101(c, err.Error())
			return
		}
		ReturnSuccess(c, "success")
		return
	}
	//判断银行卡是否存在
	if action == "BankIfExist" {
		B := model.Bank{WorkerId: int(whoMap["ID"].(uint))}
		ReturnSuccessData(c, B.BankIsExist(mysql.DB), "success")
		return
	}

	if action == "GetBackList" {
		b := make([]model.Bank, 0)
		mysql.DB.Where("worker_id=?", int(whoMap["ID"].(uint))).Find(&b)
		ReturnSuccessData(c, b, "success")
		return
	}
	if action == "GetBackListUpdate" {
		id := c.Query("id")
		name := c.Query("name")
		mail := c.Query("mail")
		phone := c.Query("phone")
		address := c.Query("address")
		if address == "" {
			ReturnErr101(c, "The collection address must not be empty")
			return
		}
		b := model.Bank{Address: address, Name: name, Phone: phone, Mail: mail, WorkerId: int(whoMap["ID"].(uint))}
		mysql.DB.Model(&model.Bank{}).Where("id=?", id).Update(&b)
		ReturnSuccess(c, "success")
		return
	}

	if action == "GetBackListDel" {
		id := c.Query("id")
		mysql.DB.Model(&model.Bank{}).Where("id=?", id).Delete(&model.Bank{})
		ReturnSuccess(c, "success")
		return
	}

	if action == "lanternSlide" {
		ls := make([]model.LanternSlide, 0)
		mysql.DB.Where("status=? and language =? ", 1, c.Query("language")).Find(&ls)
		ReturnSuccessData(c, ls, "ok")
		return
	}

	// 获取用户的最新余额
	if action == "getMoney" {
		rrr := make(map[string]interface{})
		rrr["money"] = whoMap["Balance"]
		ReturnSuccessData(c, rrr, "ok")
		return
	}

}

//日结报表
func DailyReport(c *gin.Context) {
	who, _ := c.Get("who")
	whoMap := who.(map[string]interface{})
	type ReturnData struct {
		TaskNum        int     `json:"task_num"`        //任务数量
		TaskEarnings   float64 `json:"task_earnings"`   //任务收益
		JuniorTaskNum  int     `json:"junior_task_num"` //下级任务数量
		JuniorEarnings float64 `json:"junior_earnings"` //下级任务收益
		Date           string  `json:"date"`
	}

	ce := make([]ReturnData, 0)
	//任务收益(自己) 任务数量(自己)
	mysql.DB.Raw("SELECT  SUM(tasks.price) as task_earnings ,count(*) as task_num ,date FROM task_orders   LEFT JOIN tasks  on  tasks.id=task_orders.task_id  where task_orders.status=3  and task_orders.worker_id=? GROUP BY task_orders.date  ", whoMap["ID"]).Scan(&ce)
	ce2 := make([]ReturnData, 0)
	mysql.DB.Raw("SELECT   SUM(tasks.price) as  junior_earnings,count(*) as junior_task_num,date   FROM task_orders   LEFT JOIN tasks  on  tasks.id=task_orders.task_id LEFT JOIN workers ON  workers.id=task_orders.worker_id WHERE workers.superior_id=20    GROUP BY task_orders.date  ").Scan(&ce2)
	ce3 := make([]ReturnData, 0)
	var p []string
	for _, i2 := range ce {
		p = append(p, i2.Date)
	}
	for k, i2 := range ce {
		//判断这个时间点是否已经创建了
		for k1, i4 := range ce2 {
			if i2.Date == i4.Date {
				ce[k].JuniorTaskNum = i4.JuniorTaskNum
				ce[k].JuniorEarnings = i4.JuniorEarnings
			} else {
				if tools.InArray(p, i4.Date) == false {
					ce3 = append(ce3, ce2[k1])
					p = append(p, i4.Date)
				}
			}
		}
		ce3 = append(ce3, ce[k])
	}

	ReturnSuccessData(c, ce3, "OK")
	return

}

//购买云管家
func GetCloudHousekeeper(c *gin.Context) {
	//先判断用户是否有资格购买
	who, _ := c.Get("who")
	whoMap := who.(map[string]interface{})
	config := model.Config{}
	err := mysql.DB.Where("id=?", 1).First(&config).Error
	if err != nil {
		ReturnErr101(c, "system error")
		return
	}
	if config.OpenCloudHousekeeperLevelId > whoMap["VipId"].(int) {
		vip := model.Vip{ID: uint(config.OpenCloudHousekeeperLevelId)}
		ReturnErr101(c, "Sorry, the lowest open level is "+vip.GetLevelName(mysql.DB))
		return
	}
	//判断是否已经购买了   如果已经购买 就续费
	wb := model.WorkerBalance{ID: int(whoMap["ID"].(uint)), AddBalance: 999, Kinds: 11}
	_, err = wb.AddBalanceFuc(mysql.DB)
	if err != nil {
		ReturnErr101(c, err.Error())
		return
	}
	ReturnSuccess(c, "success")
	return
}



