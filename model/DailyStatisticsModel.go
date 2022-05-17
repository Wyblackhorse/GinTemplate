/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/wangyi/GinTemplate/tools"
	"time"
)

type DailyStatistics struct {
	ID                  uint    `gorm:"primaryKey;comment:'主键'"`
	TodayRegister       int     //每日注册人数
	TodayAddVipNums     int     //今日新增vip数量
	TodayAddSubmitTask  int     //今日提交任务数量
	TodayRechargePeople int     //今日充值人数
	TodayWithdrawPeople int     //今日提现人数
	TodayRechargeMoney  float64 //今日充值
	TodayWithdrawMoney  float64 //今日提现总额
	VipYe               float64 `json:"vip_ye"` //会员余额
	Date                string  //日期
	Updated             int64
	Month               int //月
	Week                int //周

	WorkerId int `gorm:"-"`
}

func CheckIsExistModelDailyStatistics(db *gorm.DB) {
	if db.HasTable(&DailyStatistics{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&DailyStatistics{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&DailyStatistics{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}

//设置每日数据
func (d *DailyStatistics) SetEverydayData(db *gorm.DB) {
	d.Date = time.Now().Format("2006-01-02")
	d.Updated = time.Now().Unix()
	every := DailyStatistics{}
	d.Week = tools.ReturnTheWeek()
	d.Month = tools.ReturnTheMonth()
	err := db.Where("date=?", d.Date).First(&every).Error
	// select * form 表名字  where  date = 2022-05-16
	if err != nil {
		db.Save(&d)
	}
	//更新你注册
	if d.TodayRegister == 1 {
		d.TodayRegister = d.TodayRegister + every.TodayRegister
	}
	//今日新增vip数量
	if d.TodayAddVipNums == 1 {
		d.TodayAddVipNums = d.TodayAddVipNums + every.TodayAddVipNums
	}
	//今日提交任务人数的
	if d.TodayAddSubmitTask == 1 {
		d.TodayAddSubmitTask = d.TodayAddSubmitTask + every.TodayAddSubmitTask
	}

	//今日充值 	//今日充值人数
	if d.TodayRechargeMoney != 0 {
		err := db.Where("worker_id=?", d.WorkerId).Where("kinds=?", 1).Where("date=?", time.Now().Format("2006-01-02")).First(&Record{}).Error
		if err != nil {
			d.TodayRechargePeople = d.TodayRechargePeople + every.TodayRechargePeople
		}
		d.TodayRechargeMoney = d.TodayRechargeMoney + every.TodayRechargeMoney
	}

	//今日提现总额  	//今日提现人数
	if d.TodayWithdrawMoney != 0 {
		err := db.Where("worker_id=?", d.WorkerId).Where("kinds=?", 2).Where("date=?", time.Now().Format("2006-01-02")).First(&Record{}).Error
		if err != nil {
			d.TodayRechargePeople = d.TodayRechargePeople + every.TodayRechargePeople
		}
		d.TodayWithdrawMoney = d.TodayWithdrawMoney + every.TodayWithdrawMoney
	}
	//会员的余额 SELECT  SUM(balance)  FROM  workers
	db.Raw("SELECT  SUM(balance) as vip_ye  FROM  workers").Scan(&d)
	db.Model(&DailyStatistics{}).Where("id=?", every.ID).Update(&d)
}
