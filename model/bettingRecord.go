package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type BettingRecord struct {
	UserID              int     //用户ID
	UserName            string  //用户名
	EventID             int     //赛事ID
	State               string  //状态
	BreakEven           string  //保本
	BetAmount           float64 //下注金额
	Payout              float64 //派彩
	ProfitAndLoss       string  //盈/亏
	OrderID             int     //orderID
	FullHalf            string  //全/半场
	PositiveWaveCounter string  //正/反波
	Score               string  //比分
	Alliance            string  //联盟
	TopVSBottom         string  //主：客队
	Yield               float64 //收益率%
	StringOfNo          string  //串号
	TheFinalResult      string  //最终结果
	TheStartTime        string  //开赛时间
	BetOnTime           string  //下注时间
	Poundage            float64 //手续费
	PoundagePer         float64 //手续费百分号
	BetDate             string  //下注日期
}

func CheckIsExistModelBettingRecord(db *gorm.DB) {
	if db.HasTable(&BettingRecord{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&BettingRecord{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&BettingRecord{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}
