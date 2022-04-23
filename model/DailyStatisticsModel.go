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

type DailyStatistics struct {
	ID       uint   `gorm:"primaryKey;comment:'主键'"`
	Register int    //每日注册人数
	Date     string //日期
	Updated  int64
}

func CheckIsExistModelDailyStatistics(db *gorm.DB) {
	if db.HasTable(&DailyStatistics{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&DailyStatistics{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&DailyStatistics{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}

//设置每日数据
func (d *DailyStatistics) SetEverydayData(db *gorm.DB) {
	d.Date = time.Now().Format("2006-01-02")
	d.Updated = time.Now().Unix()
	every := DailyStatistics{}
	err := db.Where("date=?", d.Date).First(&every).Error
	if err != nil {
		db.Save(&d)
	}
	//更细你注册
	if d.Register == 1 {
		d.Register = d.Register + every.Register
		db.Model(&DailyStatistics{}).Where("id=?", every.ID).Update(&d)
	}

}
