package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type WalletRecord struct {
	ID             uint    `gorm:"primaryKey"`
	SerialNumber   int     //编号
	UserName       string  //用户名
	Amount         float64 //金额
	Kinds          string  //类型
	BeforeTheValue int     //前值
	AfterTheValue  int     //后值
	Serial         string     //串号
	Note           string  //备注
	Created        string  //创建时间
	Updated        string  //更新时间
}
func CheckIsExistModelWalletRecord(db *gorm.DB) {
	if db.HasTable(&WalletRecord{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&WalletRecord{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&WalletRecord{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}