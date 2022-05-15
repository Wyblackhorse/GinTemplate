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

type Vip struct {
	ID        uint `gorm:"primaryKey;comment:'主键'"`
	Name      string
	Level     int     //等级
	Price     float64 `gorm:"type:decimal(10,2)"` //价格
	Account   float64 `gorm:"type:decimal(10,2)"` //每单的收入
	TaskTimes int     //任务
	Created   int64   //创建时间
}

func CheckIsExistModelVip(db *gorm.DB) {
	if db.HasTable(&Vip{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Vip{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Vip{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
			m := Vip{ID: 1, Level: 1, Price: 0, Account: 2, TaskTimes: 5, Name: "Estuary"}
			m.Add(db)
			m = Vip{ID: 2, Level: 1, Price: 960, Account: 8, TaskTimes: 5, Name: "VIP1"}
			m.Add(db)
			m = Vip{ID: 3, Level: 2, Price: 2240, Account: 8, TaskTimes: 12, Name: "VIP2"}
			m.Add(db)
			m = Vip{ID: 4, Level: 3, Price: 5200, Account: 8, TaskTimes: 30, Name: "VIP3"}
			m.Add(db)
			m = Vip{ID: 5, Level: 4, Price: 16000, Account: 12, TaskTimes: 62, Name: "VIP4"}
			m.Add(db)
			m = Vip{ID: 6, Level: 5, Price: 32800, Account: 12, TaskTimes: 130, Name: "VIP5"}
			m.Add(db)
			m = Vip{ID: 7, Level: 6, Price: 68000, Account: 12, TaskTimes: 270, Name: "VIP6"}
			m.Add(db)
			m = Vip{ID: 8, Level: 7, Price: 144000, Account: 16, TaskTimes: 470, Name: "VIP7"}
			m.Add(db)
			m = Vip{ID: 9, Level: 8, Price: 216000, Account: 20, TaskTimes: 600, Name: "VIP8"}
			m.Add(db)
			m = Vip{ID: 10, Level: 9, Price: 328000, Account: 20, TaskTimes: 960, Name: "VIP9"}
			m.Add(db)

		}
	}
}

//添加会员
func (v *Vip) Add(db *gorm.DB) {
	v.Created = time.Now().Unix()
	err := db.Where("name=?", v.Name).First(&Vip{}).Error
	if err != nil {
		db.Save(&v)
	}
}

//获取会员等级
func (v *Vip) GetLevelName(db *gorm.DB) string {
	vip := Vip{}
	err := db.Where("id=?", v.ID).First(&vip).Error
	if err != nil {
		return ""
	}
	return vip.Name
}

//判断会员是否存在 如果存在返回价格
func (v *Vip) ReturnVipPrice(db *gorm.DB) (bool, float64) {
	//判断vip 是否存在
	vip := Vip{}
	err := db.Where("id=?", v.ID).First(&vip).Error
	if err != nil {
		return false, 0
	}

	return true, vip.Price
}
