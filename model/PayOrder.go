package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

// PayOrder 支付订单
type PayOrder struct {
	ID          uint    `gorm:"primaryKey;comment:'主键'"`
	TxHash      string  //转账hash 值
	BlockNumber int     //区块号
	Timestamp   int64   //时间戳
	FromAddress string  //转账地址
	ToAddress   string  //收账地址
	Amount      float64 `gorm:"type:decimal(10,2)"` //金额
	Token       string  //token
	UserID      string  //用户id
	Created     int64
	Date        string
}

func CheckIsExistModePayOrder(db *gorm.DB) {
	if db.HasTable(&PayOrder{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&PayOrder{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&PayOrder{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}

// IfIsExitsThisData 判断数据是否存在
func (p *PayOrder) IfIsExitsThisData(db *gorm.DB) bool {
	err := db.Where("tx_hash=?", p.TxHash).First(&PayOrder{}).Error
	if err != nil {
		return false
	}
	return true
}
