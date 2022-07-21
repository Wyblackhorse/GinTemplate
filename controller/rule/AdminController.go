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
	"github.com/wangyi/GinTemplate/dao/redis"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"time"
)

//登录
func Login(c *gin.Context) {
	var admin LoginRequest
	if err := c.ShouldBind(&admin); err != nil {
		tools.JsonWrite(c, -101, []string{}, "err:"+err.Error())
		return
	}
	adminTwo := model.AdminModel{}
	err := mysql.DB.Where("username=?", admin.Username).Where("password=?", admin.Password).First(&adminTwo).Error
	if err != nil {
		tools.JsonWrite(c, -101, []string{}, "账户或者密码错误!")
		return
	}
	//查看账号状态
	if adminTwo.Status == 2 {
		tools.JsonWrite(c, -101, []string{}, "该账号已经停用")
		return
	}
	//if adminTwo.GoogleCode == "" {
	//	secret, _, qrCodeUrl := tools.InitAuth(adminTwo.Username)
	//	err := mysql.DB.Model(&model.AdminModel{}).Where("id=?", adminTwo.ID).Update(model.AdminModel{GoogleCode: secret}).Error
	//	if err != nil {
	//		tools.JsonWrite(c, -101, []string{}, "err:"+err.Error())
	//		return
	//	}
	//	tools.JsonWrite(c, -102, map[string]string{"codeUrl": qrCodeUrl}, "请先绑定谷歌账号")
	//	return
	//}
	////校验谷歌验证
	//verifyCode, _ := tools.NewGoogleAuth().VerifyCode(adminTwo.GoogleCode, admin.GoogleCode)
	//if !verifyCode {
	//	tools.JsonWrite(c, -101, map[string]string{}, "谷歌验证失败")
	//	return
	//}
	//登录成功
	redis.Rdb.Set("Admin_Login_Token_"+adminTwo.Token, adminTwo.Username, time.Second*3600*24)
	//登录日志
	log := model.Loggers{Content: "登陆成功", WriteId: int(adminTwo.ID)}
	log.AdminAddErrorLoggers(mysql.DB)
	tools.JsonWrite(c, 200, adminTwo.Token, "登录成功")
	return
}

//云管家自动去做任务
func CloudHousekeeperDoTask(c *gin.Context) {
	ws := make([]model.Worker, 0)
	mysql.DB.Where("cloud_housekeeper=? and cloud_housekeeper_expire > ?", 2, time.Now().Unix()).Find(&ws)
	fmt.Println(len(ws))
	for _, i2 := range ws {
		m := 10.0
		//判断用户的vip等级  可以做的任务量
		if i2.VipExpire > time.Now().Unix() {
			vips := model.Vip{}
			err := mysql.DB.Where("id=?", i2.VipId).First(&vips).Error
			if err == nil {
				//可以做的任务量
				//vips.TaskTimes    // price
				var  count  int
				mysql.DB.Model(&model.TaskOrder{}).Where("status=?  and  worker_id=? and   date=?", 3, i2.ID, time.Now().Format("2006-01-02")).Count(&count)
				m = float64(vips.TaskTimes-count) * vips.Account
				if m>0 {
					go RobotTask(vips.TaskTimes, vips.Account, int(i2.ID))

				}
			}
		} else {
			//默认你是实习生   每天只可以做5单   每单两块钱
			go RobotTask(5, 2, int(i2.ID))
		}



		ws := model.WorkerBalance{ID: int(i2.ID), AddBalance: m, Kinds: 3}
		ws.AddBalanceFuc(mysql.DB)
	}
	//对用户进行余额操作
	ReturnSuccess(c, "success")

}

//机器人生成任务
func RobotTask(TaskTime int, price float64, workerId int) {






	mysql.DB.Where("")
	for i := 0; i < TaskTime; i++ {
		type TaskOrder struct {
			ID       uint   `gorm:"primaryKey;comment:'主键'"`
			TaskId   int    //任务 id
			WorkerId int    //玩家id
			Status   int    // 状态1 进行中  2审核中 3已完成 4以失败 5恶意 6已放弃
			Created  int64  //创建时间
			Updated  int64  //更新时间
			ImageUrl string //图片地址
			Date     string //日期
			Month    int    //月
			Week     int    //周
		}

		TAS := model.TaskOrder{TaskId: 0, WorkerId: workerId, Status: 3, Created: time.Now().Unix(), Updated: time.Now().Unix(), ImageUrl: "", Date: time.Now().Format("2006-01-02"), Month: tools.ReturnTheMonth(), Week: tools.ReturnTheWeek(), Robot: 2}
		mysql.DB.Save(&TAS)
	}

}
