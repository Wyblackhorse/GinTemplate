package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Withdraw struct {
	ID                       uint    `gorm:"primaryKey"`
	RecordId                 int     //记录ID
	TheUserId                int     //用户ID
	TheUserName              string  //用户名
	WithdrawalAmount         float64 //提现金额
	Status                   string  //状态
	Kinds                    int     //类型
	TheActualAmount          float64 //实际金额
	TheSettlementAmount      float64 //结算金额
	Poundage                 float64 //手续费
	Rate                     float64 //税率
	OrderNo                  string  // 单号
	Classification           string  //分类
	ChannelID                int     //渠道ID
	ThirdPartyTrackingNumber string  //第三方单号
	Date                     string  //提现时间
}

func CheckIsExistModelWithdraw(db *gorm.DB) {
	if db.HasTable(&Withdraw{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Withdraw{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Withdraw{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}
