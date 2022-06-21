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

type Collection struct {
	ID         uint    `gorm:"primaryKey;comment:'主键'"`
	TaskUrl    string  //采集的地址
	Kinds      int     //应用id
	Status     int     //1  没有采集  2采集完毕  3正在采集
	TaskType   int     //1 点赞  2转发
	Expiry     int     //有效期
	TaskNum    int     //任务数量
	Price      float64 //任务单价
	DemandSide string  //需求方
	Created    int64
	Remark     string //备注
	TaskLevel  int    //任务等级

}

func CheckIsExistModelCollection(db *gorm.DB) {
	if db.HasTable(&Collection{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Collection{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Collection{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
			//初始化 数据
		}
	}
}

//添加
func (c *Collection) Add(db *gorm.DB) (bool, error) {
	err := db.Where("task_url=?", c.TaskUrl).First(&Collection{}).Error
	if err == nil {
		return false, eeor.OtherError("不要重复添加")
	}
	c.Created = time.Now().Unix()
	c.Status = 1
	err = db.Save(&c).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
