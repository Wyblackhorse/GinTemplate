/**
 * @Author $
 * @Description //TODO $幻灯片
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

type LanternSlide struct {
	ID         uint   `gorm:"primaryKey;comment:'主键'"`
	UrlAddress string //图片地址
	Status     int    //状态    1正常  2禁用
	Language   string //所属语言
	Created    int64
}

func CheckIsExistModelLanternSlide(db *gorm.DB) {
	if db.HasTable(&LanternSlide{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&LanternSlide{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&LanternSlide{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}

//添加
func (l *LanternSlide) Add(db *gorm.DB) bool {
	l.Created = time.Now().Unix()
	l.Status = 1
	err := db.Save(&l).Error
	if err != nil {
		return false
	}
	return true
}
