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
	"time"
)

type Role struct {
	ID      uint   `gorm:"primaryKey;comment:'主键'"`
	Name    string //角色名字
	Status  int    // 1 正常 2封禁
	Created int64  //创建时间
}

func CheckIsExistModelRole(db *gorm.DB) {
	if db.HasTable(&Role{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Role{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Role{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
			//初始化 数据
			role := Role{ID: 1, Name: "超级管理员", Status: 1, Created: time.Now().Unix()}
			db.Save(&role)
		}
	}
}
