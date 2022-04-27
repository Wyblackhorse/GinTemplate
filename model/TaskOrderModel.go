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

type TaskOrder struct {
	ID         uint   `gorm:"primaryKey;comment:'主键'"`
	TaskId     int    //任务 id
	WorkerId   int    //玩家id
	Status     int    // 状态1 进行中2审核中 3已完成 4以失败 5恶意 6已放弃
	Created    int64  //创建时间
	Updated    int64  //更新时间
	ImageUrl   string //图片地址
	Date       string //日期
	TaskName   string `gorm:"-"`
	WorkerName string `gorm:"-"`
}

func CheckIsExistModelTaskOrder(db *gorm.DB) {
	if db.HasTable(&TaskOrder{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&TaskOrder{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&TaskOrder{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}
