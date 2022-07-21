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
	"sync"
	"time"
)

type Task struct {
	ID           uint    `gorm:"primaryKey;comment:'主键'"`
	ApplyId      int     //应用id
	ApplyName    string  `gorm:"-"`
	Remark       string  //备注(标题)
	TaskType     int     //任务类型   1点赞  2转发  3点赞转发
	TaskUrl      string  //任务地址
	EndTime      int64   //任务结束时间
	Price        float64 `gorm:"type:decimal(10,2)"` //价格
	Status       int     //任务状态 1正常 2结束 3取消s
	DemandSide   string  //需求方
	TaskLevel    int     //任务级别  1
	TaskNum      int     //任务总数量
	Created      int64   //创建时间
	AlreadyGet   int     //  已经领取的名额
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

type GetTaskData struct {
	GetTaskDataLock sync.RWMutex //锁
	WorkerVipId     int          //用户的vipId
	WorkerId        int          //用户id
	ApplyId         int          //应用id
	TaskId          int          //任务id
	TaskNum         int          //任务数量
	Status          int          // 6  撤销任务
	CreditScore     int          //信用分

}

//领取任务
func (g *GetTaskData) GetTask(db *gorm.DB) (bool, error) {

	if g.Status == 6 {
		g.GetTaskDataLock.Lock() //上写锁
		//解除写锁
		defer g.GetTaskDataLock.Unlock()
		db.Model(&Task{}).Where("id=?", g.TaskId).UpdateColumn("already_get", gorm.Expr("already_get+ ?", 1))
		return true, nil
	}

	//获取用户的 可以做的订单可以领取的订单数量
	vip := Vip{}
	err := db.Where("id=?", g.WorkerVipId).First(&vip).Error
	if err != nil {
		return false, err
	}
	//查询今天已经领取获取这完成的 订单数量
	var haveDone int
	db.Model(&TaskOrder{}).Where("worker_id=?", g.WorkerId).Where("status=?  or status=?  or status=?", 1, 2, 3).Where("date=?", time.Now().Format("2006-01-02")).Count(&haveDone)
	//信用分小于30分的时候  任务减半 	//今日的次数已经用完了
	if g.CreditScore < 31 {
		if g.CreditScore <= 0 {
			return false, eeor.OtherError("Sorry, the account is illegal")
		}else {
			if haveDone >= (vip.TaskTimes / 2) {
				return false, eeor.OtherError("I have run out of times today")
			}
		}
	} else {
		if haveDone >= vip.TaskTimes {
			return false, eeor.OtherError("I have run out of times today")
		}
	}

	//领取
	//上读锁
	g.GetTaskDataLock.RLock()
	//获取今日还剩余的任务数量
	task := Task{}
	err = db.Where("id=? and status=?", g.TaskId, 1).First(&task).Error
	if err != nil {
		g.GetTaskDataLock.RUnlock() //解除读锁
		return false, eeor.OtherError("Mission does not exist")
	}
	if task.AlreadyGet <= 0 {
		g.GetTaskDataLock.RUnlock() //解除读锁
		return false, eeor.OtherError("Sorry, this assignment has been taken")
	}
	//查看任务是否过期
	if task.EndTime < time.Now().Unix() {
		g.GetTaskDataLock.RUnlock() //解除读锁
		return false, eeor.OtherError("Sorry, The task has expired")
	}
	db = db.Begin()
	//解除读锁
	g.GetTaskDataLock.RUnlock()
	//上写锁
	g.GetTaskDataLock.Lock()
	//写锁解除
	defer g.GetTaskDataLock.Unlock()

	//生成任务
	t := TaskOrder{WorkerId: g.WorkerId, TaskId: g.TaskId}
	_, err = t.SetTaskOrder(db)
	if err != nil {
		return false, err
	}
	//总任务数量 -1
	err = db.Model(&Task{}).Where("id=?", task.ID).Update(map[string]interface{}{"already_get": task.AlreadyGet - 1}).Error
	if err != nil {
		db.Rollback() //事务回滚
		return false, err
	}
	db.Commit()
	return true, nil
}

//任务是否存在
func (t *Task) IsExist(db *gorm.DB) bool {
	err := db.Where("id=?", t.ID).First(&t).Error
	if err != nil {
		return false
	}
	return true
}
