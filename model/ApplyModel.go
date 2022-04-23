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

// 应用

type Apply struct {
	ID      uint   `gorm:"primaryKey;comment:'主键'"`
	Name    string //菜单名称
	Status  int    //1正常 2禁用
	Created int64  //创建时间
}

//添加数据
func (a *Apply) AddApply(db *gorm.DB) int {
	a.Status = 1
	a.Created = time.Now().Unix()
	err := db.Where("name=?", a.Name).First(&Apply{}).Error
	if err == nil {
		//不要重复添加
		return -1
	}
	db.Save(&a)
	return 1
}

func CheckIsExistModelApply(db *gorm.DB) {
	if db.HasTable(&Apply{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Apply{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Apply{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
			a := Apply{ID: 1, Name: "YouTube"}
			a.AddApply(db)
			a = Apply{ID: 2, Name: "line"}
			a.AddApply(db)
			a = Apply{ID: 3, Name: "Facebook"}
			a.AddApply(db)
			a = Apply{ID: 4, Name: "TikTok"}
			a.AddApply(db)
			a = Apply{ID: 5, Name: "Instagram"}
			a.AddApply(db)
		}
	}
}




