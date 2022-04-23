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

type Config struct {
	ID                 uint `gorm:"primaryKey;comment:'主键'"`
	NeedInvitationCode int  `gorm:"int(10);default:2"`  //是否需要邀请码  1需要2不需要
	CreditScore        int  `gorm:"int(10);default:60"` //设置默认信用分
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
			db.Save(&Config{ID: 1})
		}
	}
}
