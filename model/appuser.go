package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type AppUser struct {
	ID           uint    `gorm:"primaryKey"`

	UserNumber         int    //用户序号
	TheHigherTheID     int    //上级ID
	UpperLayerUserName string //上级用户名
	GeneralAgentID     int    //总代理ID
	TheGeneralAgentOf  string //总代理名
	Username        string //用户名
	MobilePhoneNo      string    //手机号
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