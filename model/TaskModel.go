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

type Task struct {
	ID           uint    `gorm:"primaryKey;comment:'主键'"`
	ApplyId      int     //应用id
	ApplyName    string  `gorm:"-"`
	Remark       string  //备注
	TaskType     int     //任务类型   1点赞  2转发
	TaskUrl      string  //任务地址
	EndTime      int64   //任务结束时间
	Price        float64 `gorm:"type:decimal(10,2)"` //价格
	Status       int     //任务状态 1正常 2结束 3取消s
	TaskLevel    int     //任务级别  1
	TaskNum      int     //任务数量
	Created      int64   //创建时间
	WorkerStatus int     `gorm:"-"` //1 没有做 2已结提交了图片

}

func CheckIsExistModelTask(db *gorm.DB) {
	if db.HasTable(&Task{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Task{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Task{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}
