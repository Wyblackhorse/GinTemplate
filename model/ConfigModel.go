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
	NeedInvitationCode          int     `gorm:"int(10);default:2"`                  //是否需要邀请码  1需要2不需要
	CreditScore                 int     `gorm:"int(10);default:60"`                 //注册设置默认信用分
	DoneFirstTask               int     `gorm:"int(10);default:1"`                  //完成当日首次任务
	OverTimeForTask             int     `gorm:"int(10);default:1"`                  //没有在规定的时间内完成任务
	NoRequireToTask             int     `gorm:"int(10);default:1"`                  //没有按照要求完成任务
	InviteRewards               float64 `gorm:"type:decimal(10,2);default:0"`       //邀请奖励
	SuperiorBackPer             float64 `gorm:"type:decimal(10,2);default:0"`       //上级返点
	NextSuperiorBackPer         float64 `gorm:"type:decimal(10,2);default:0"`       //上上级返点
	NextNextSuperiorBackPer     float64 `gorm:"type:decimal(10,2);default:0"`       //上上上级返点
	LowWithdrawal               float64 `gorm:"type:decimal(10,2);default:0"`       //最低提现金额
	HighWithdrawal              float64 `gorm:"type:decimal(10,2);default:9999999"` //最高提现金额
	WithdrawalCharge            float64 `gorm:"type:decimal(10,2);default:0"`       //提现手续费
	ReferralBonusesTimes        int     `gorm:"default: 1"`                         //推荐奖励次数     1不重复发放  2重复发放
	ForegroundLanguage          string  `gorm:"default: '简体中文'"`                    //前台默认  语言
	LinkOfTheService            string  //客服地址
	TelegramAddress             string  //飞机地址
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
			db.Save(&Config{ID: 1, LowWithdrawal: 10})
		}
	}
}
