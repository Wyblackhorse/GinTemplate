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
	eeor "github.com/wangyi/GinTemplate/error"
	"time"
)

type Bank struct {
	ID       uint `gorm:"primaryKey;comment:'主键'"`
	WorkerId int
	Name     string
	Address  string
	Phone    string
	Mail     string
	Created  int64
	WorkerName  string  `gorm:"-"`
}

func CheckIsExistModelBank(db *gorm.DB) {
	if db.HasTable(&Bank{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Bank{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Bank{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}

//添加  银行卡地址
func (b *Bank) Add(db *gorm.DB) (bool, error) {
	//判断这个条数据是否已经存在
	err := db.Where("worker_id=?", b.WorkerId).Where("address=?", b.Address).First(&b).Error
	if err == nil {
		return false, eeor.OtherError("Don't add more than once")
	}
	b.Created = time.Now().Unix()
	db.Save(&b)
	return true, nil
}

func (b *Bank) BankIsExist(db *gorm.DB) Bank {
	err := db.Where("worker_id=?", b.WorkerId).First(&b).Error
	if err != nil {
		return Bank{}
	}
	return *b
}
