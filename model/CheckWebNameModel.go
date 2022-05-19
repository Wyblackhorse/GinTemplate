package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	eeor "github.com/wangyi/GinTemplate/error"
	"time"
)

// CheckWebName WebName 域名网站
type CheckWebName struct {
	ID        uint `gorm:"primaryKey;comment:'主键'"`
	WebNameId int  //域名id
	MatchUrl  string
	Status    int //状态  1 无效  2有效 3未检测
	Created   int64
	Updated   int64
	WebName  string  `gorm:"-"`
}

func CheckIsExistModelCheckWebName(db *gorm.DB) {
	if db.HasTable(&CheckWebName{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&CheckWebName{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&CheckWebName{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}

func (c *CheckWebName) AddCheckWebName(db *gorm.DB) (bool, error) {
	//判断是否存在
	err := db.Where("web_name_id=?", c.WebNameId).Where("match_url=?", c.MatchUrl).First(&CheckWebName{}).Error
	if err == nil {
		return false, eeor.OtherError("不要重复添加")
	}
	c.Created = time.Now().Unix()
	c.Status = 3
	db.Save(&c)
	return true, nil
}
