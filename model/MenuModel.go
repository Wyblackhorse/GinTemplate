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

type Menu struct {
	ID uint `gorm:"primaryKey;comment:'主键'"`
}




func CheckIsExistModelMenu(db *gorm.DB) {
	if db.HasTable(&Menu{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Menu{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Menu{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}
