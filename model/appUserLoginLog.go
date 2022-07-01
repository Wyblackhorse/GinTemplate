package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type AppUserLoginLog struct {
	Create               string  //创建时间
	Update               string //更新时间
	Remark                 string //备注
	SerialNumber         int    //序号
	UserID               int    //用户编号
	UserAccount          string //用户账号
	LoginID              string //登录状态
	TheLoginAddress      string //登录地址
	TheLoginSite         string //登录地点
	Browser              string //浏览器
	OperatingSystem      string //操作系统
	OperatingInformation string //操作信息
	AccessTime           string //访问时间
}





func CheckIsExistModelAppUserLoginLog(db *gorm.DB) {
	if db.HasTable(&AppUserLoginLog{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&AppUserLoginLog{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&AppUserLoginLog{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}