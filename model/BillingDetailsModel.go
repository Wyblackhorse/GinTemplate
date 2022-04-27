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

type BillingDetails struct {
	ID          uint `gorm:"primaryKey;comment:'主键'"`
	WorkerId    int
	ChangeMoney float64 `gorm:"type:decimal(10,2)"`
	InitMoney   float64 `gorm:"type:decimal(10,2)"`
	NowMoney    float64 `gorm:"type:decimal(10,2)"`
	Kinds       int     //类型 1充值  2提现 3做单任务 4购买业务 5佣金奖励
	Created     int64
	WorkerName  string `gorm:"-"`
}



func CheckIsExistModelBillingDetails(db *gorm.DB) {
	if db.HasTable(&BillingDetails{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&BillingDetails{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&BillingDetails{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}
