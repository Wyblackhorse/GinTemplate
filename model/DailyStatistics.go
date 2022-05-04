package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

// DailyStatistics 每日统计
type DailyStatistics struct {
	ID              uint    `gorm:"primaryKey;comment:'主键'"`
	RechargeAccount float64 `gorm:"type:decimal(10,2)"` //充值的金额
	RechargeNums    int     //充值的笔数
	Created         int64   // 创建时间
	Updated         int64   //更新时间
	Date            string  //日期
}

func CheckIsExistModeDailyStatistics(db *gorm.DB) {
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

func (d *DailyStatistics) UpdateDailyStatistics(db *gorm.DB) {
	//判断今日的数据是否存在
	da := DailyStatistics{}
	err := db.Where("date=?", time.Now().Format("2006-01-02")).First(&da).Error
	if err != nil {
		//不存在这条数据 创建
		db.Save(&DailyStatistics{Date: time.Now().Format("2006-01-02"), Created: time.Now().Unix(), RechargeNums: 1, RechargeAccount: d.RechargeAccount})
		return
	}
	//更新
	new := DailyStatistics{}
	new.RechargeNums = da.RechargeNums + 1
	new.Updated = time.Now().Unix()
	new.RechargeAccount = da.RechargeAccount + d.RechargeAccount
	db.Model(DailyStatistics{}).Where("id=?", da.ID).Update(&new)
}
