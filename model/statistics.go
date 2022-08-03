package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Statistics struct {
	ID                     uint    `gorm:"primaryKey"`
	InAllWithdrawNoBetNums int     `json:"in_all_withdraw_no_bet_nums"` //充值没有下注的人数
	InAllWithdrawBetNums   int     `json:"in_all_withdraw_bet_nums"`    //充值下注活跃人数
	InAllWithdrawNums      int     `json:"in_all_withdraw_nums"`        //充值人数
	InAllWithdrawFirst     int     `json:"in_all_withdraw_first"`       //首冲人数
	InAllPrivatelyBetNums  int     `json:"in_all_privately_bet_nums"`   //私自下注人数  下注订单并且下注不是保本的
	InAllWithdrawDeposit   float64 `json:"in_all_withdraw_deposit"`     //存款提现数据
	InAllCommission        float64 `json:"in_all_commission"`           //佣金数据

	TodayWithdrawNoBetNums int     `json:"today_withdraw_no_bet_nums"` //今日充值没有下注的人数
	TodayWithdrawBetNums   int     `json:"today_withdraw_bet_nums"`    //今日充值下注活跃人数
	TodayWithdrawNums      int     `json:"today_withdraw_nums"`        //今日充值人数
	TodayWithdrawFirst     int     `json:"today_withdraw_first"`       //今日首冲人数
	TodayPrivatelyBetNums  int     `json:"today_privately_bet_nums"`   //今日私自下注人数  下注订单并且下注不是保本的
	TodayAllBetNums        int     `json:"today_all_bet_nums"`         //今日总下注人数
	TodayWithdrawDeposit   float64 `json:"today_withdraw_deposit"`     //今日存款提现数据
	TodayCommission        float64 `json:"today_commission"`           //今日佣金数据
	Created                int64   `json:"created"`
	Updated                int64   `json:"updated"`
	Date                   string  `json:"date"`
}

type User struct {
	Username          string `json:"username"`              //用户名
	TheGeneralAgentOf string `json:"upper_layer_user_name"` //顶级总代
}

type Detail struct {
	TodayPrivatelyBetNumsDetail  []User //私自下注
	TodayWithdrawNoBetNumsDetail []User //充值未投注
	TodayWithdrawFirstDetail     []User //首充会员名单
	TodayWithdrawBetNumsDetail   []User //下注人员明细
}

func CheckIsExistModelTotalWithdraw(db *gorm.DB) {
	if db.HasTable(&Statistics{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Statistics{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Statistics{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}
