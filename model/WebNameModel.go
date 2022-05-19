package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	eeor "github.com/wangyi/GinTemplate/error"
	"strings"
	"time"
)

// WebName 域名网站
type WebName struct {
	ID          uint   `gorm:"primaryKey;comment:'主键'"`
	Url         string //域名
	MatchingUrl string
	Suffix      string
	Status      int //状态  1 正常
	Created     int64
}

func CheckIsExistModelWebName(db *gorm.DB) {
	if db.HasTable(&WebName{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&WebName{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&WebName{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}

	}
}

// AddWebName 添加域名
func (w *WebName) AddWebName(db *gorm.DB) (bool, error) {
	err := db.Where("url=?", w.Url).First(&WebName{}).Error
	if err == nil {
		return false, eeor.OtherError("不要重复添加")
	}
	w.Status = 1
	w.Created = time.Now().Unix()
	w.Suffix = strings.Split(w.Url, ".")[1]
	w.MatchingUrl = strings.Split(w.Url, ".")[0]
	err = db.Save(&w).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

//判断  是否存在
