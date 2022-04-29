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
	"strconv"
	"time"
)

type Record struct {
	ID       uint `gorm:"primaryKey;comment:'主键'"`
	WorkerId int  //用户id

	Kinds     int     //类型 1充值  2提现   4购买业务 5佣金奖励 6充值到余额宝
	Money     float64 `gorm:"type:decimal(10,2)"` //购买金额
	Status    int     //1已完成  2审核中  3失败
	RecordNum string
	Created   int64
	Updated   int64
}

func CheckIsExistModelRecord(db *gorm.DB) {
	if db.HasTable(&Record{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Record{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Record{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
			//初始化 数据
		}
	}
}

//添加订单
func (r *Record) AddRecord(db *gorm.DB) (bool, error) {
	r.Created = time.Now().Unix()
	r.RecordNum = time.Now().Format("20060102") + strconv.FormatFloat(float64(time.Now().Unix()), 'f', 0, 64) + strconv.Itoa(r.WorkerId)
	err := db.Save(&r).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
