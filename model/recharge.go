package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Recharge struct {
	ID           uint    `gorm:"primaryKey"`
	UserId       int     //用户id
	Username     string  //用户名
	Recharge     float64 //充值金额
	CloseAccount float64 //结算货币金额
	TopUpAward   float64 //充值奖励
	Status       string  //状态
	Kinds        int     //类型
	Serial       string  //串号
	ThreeOrders  string  //三方单号
	TopUpChannel int     //充值渠道
	Classify     string  //分类
	Created      string  //创建时间
	Updated      string  //更新时间
	Remark       string  //备注
	Date         string  // 日期
	DateTimestamp    int64 `json:"date_timestamp"`
}

func CheckIsExistModelRecharge(db *gorm.DB) {
	if db.HasTable(&Recharge{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Recharge{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Recharge{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}
