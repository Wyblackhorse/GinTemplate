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
//系统设置
type Config struct {
	ID                          uint    `gorm:"primaryKey;comment:'主键'"`
	NeedInvitationCode          int     `gorm:"int(10);default:2"`  //是否需要邀请码  1需要2不需要
	CreditScore                 int     `gorm:"int(10);default:60"` //设置默认信用分
	InviteRewards               float64 `gorm:"type:decimal(10,2)"` //邀请奖励
	SuperiorBackPer             float64 `gorm:"type:decimal(10,2)"` //上级返点
	NextSuperiorBackPer         float64 `gorm:"type:decimal(10,2)"` //上上级返点
	NextNextSuperiorBackPer     float64 `gorm:"type:decimal(10,2)"` //上上上级返点
	LowWithdrawal               float64 `gorm:"type:decimal(10,2)"` //最低提现金额
	WithdrawalCharge            float64 `gorm:"type:decimal(10,2)"` //提现手续费
	OpenCloudHousekeeperLevelId int     //开通云管家的最低等级
}

func CheckIsExistModelConfig(db *gorm.DB) {
	if db.HasTable(&Config{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Config{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Config{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
			//初始化 数据
			db.Save(&Config{ID: 1,LowWithdrawal: 10})
		}
	}
}
