package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type AppUser struct {
	ID                 uint   `gorm:"primaryKey"`
	UserNumber         int    //用户序号
	TheHigherTheID     int    //上级ID
	UpperLayerUserName string //上级用户名
	GeneralAgentID     int    //总代理ID
	TheGeneralAgentOf  string //总代理名
	Username           string //用户名
	MobilePhoneNo      string //手机号
	UserMailbox        string //用户邮箱
	State              string //状态
	RegistrationTime   string //注册时间
	RegisteredIP       string //注册IP
	InviteCode         string //邀请码
	RealName           string //实名
	TestNo             string //是否测试号
	Grouping           string //分组
	LastLoginTime      string //最后登录时间
	Updated            string //更新时间
}

type DetailAppUser struct {
	UserNumber             int     `json:"user_number"`              //会员id
	UpperLayerUserName     string  `json:"upper_layer_user_name"`    //会员账号
	RegistrationTime       string  `json:"registration_time"`        //注册时间
	LastLoginTime          string  `json:"last_login_time"`          //最后登录时间
	SuperiorList           string  `json:"superior_list"`            //上级列表
	DirectlySubordinateNum int     `json:"directly_subordinate_num"` //直属下级个数
	TeamNums               int     `json:"team_nums"`                //团队人数
	Money                  float64 `json:"money"`                    ///用户余额
	Recharge               float64 `json:"recharge"`                 //充值
	Withdraw               float64 `json:"withdraw"`                 //提现
	BetMoney               float64 `json:"bet_money"`                //下注金额
	BetNum                 int     `json:"bet_num"`                  //下注次数
	SaveDifference         float64 `json:"save_difference"`//存提差
}

type DetailAppUserArray struct {
	Myself DetailAppUser
	One    DetailAppUser
	Two    DetailAppUser
	Three  DetailAppUser
}

func CheckIsExistModelAppUser(db *gorm.DB) {
	if db.HasTable(&AppUser{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&AppUser{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&AppUser{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}
