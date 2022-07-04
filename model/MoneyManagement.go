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
)

//理财产品
type MoneyManagement struct {
	ID       uint    `gorm:"primaryKey;comment:'主键'"`
	WorkerId int     //用户id
	YuEBaoId int     //产品id
	Money    float64 `gorm:"type:decimal(10,2)"` //购买金额
	Status   int     //1进行中  2结束
	Created  int64   //创建时间
	EndTime  int64   //结束时间
	WorkerName   string  `gorm:"-"`
	YuEBaoName   string  `gorm:"-"`
	Day          int     `gorm:"-"`
	InterestRate float64 `gorm:"-"`
}

func CheckIsExistModelMoneyManagement(db *gorm.DB) {
	if db.HasTable(&MoneyManagement{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&MoneyManagement{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&MoneyManagement{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}
