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
	"github.com/wangyi/GinTemplate/tools"
	"time"
)

type TaskOrder struct {
	ID         uint   `gorm:"primaryKey;comment:'主键'"`
	TaskId     int    //任务 id
	WorkerId   int    //玩家id
	Status     int    // 状态1 进行中  2审核中 3已完成 4以失败 5恶意 6已放弃
	Created    int64  //创建时间
	Updated    int64  //更新时间
	ImageUrl   string //图片地址
	Date       string //日期
	Month      int    //月
	Week       int    //周
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

//返回个人信息 (收益统计)
type ReturnResult struct {
	TodayGetMoney       float64 `json:"today_get_money"`      //今日收益
	YesterdayGetMoney   float64 `json:"yesterday_get_money"`  //昨日收益
	ThisWeekGetMoney    float64 `json:"this_week_get_money"`  //本周收益
	ThisMonthGetMoney   float64 `json:"this_month_get_money"` //本月收益
	LastMonthGetMoney   float64 `json:"last_month_get_money"` //上月收益
	AllGetMoney         float64 `json:"all_get_money"`        //总收益
	TodayDoneTaskNum    int     //今日完成单数
	TodayResidueTaskNum int     //今日剩余单数
	Balance             float64 //用户余额

}

//获取收益 昨日  今日 本周 本月 上个月
func (t *TaskOrder) GetEarnings(db *gorm.DB, taskTimes int) ReturnResult {

	//今日收益

	var result ReturnResult
	//今日收益
	db.Table("task_orders").
		Select("sum(tasks.price) as today_get_money").Joins("LEFT JOIN tasks ON  task_orders.task_id=tasks.id").
		Where("task_orders.worker_id=?", t.WorkerId).Where("task_orders.status=?", 3).Where("date=?", time.Now().Format("2006-01-02")).Scan(&result)
	//昨日收益
	db.Table("task_orders").
		Select("sum(tasks.price) as yesterday_get_money").Joins("LEFT JOIN tasks ON  task_orders.task_id=tasks.id").
		Where("task_orders.worker_id=?", t.WorkerId).Where("task_orders.status=?", 3).Where("date=?", time.Now().AddDate(0, 0, -1).Format("2006-01-02")).Scan(&result)

	//获取本周收益
	db.Table("task_orders").
		Select("sum(tasks.price) as this_week_get_money").Joins("LEFT JOIN tasks ON  task_orders.task_id=tasks.id").
		Where("task_orders.worker_id=?", t.WorkerId).Where("task_orders.status=?", 3).
		Where("week=?", tools.ReturnTheWeek()).Scan(&result)

	//获取本月
	db.Table("task_orders").
		Select("sum(tasks.price) as this_month_get_money").Joins("LEFT JOIN tasks ON  task_orders.task_id=tasks.id").
		Where("task_orders.worker_id=?", t.WorkerId).Where("task_orders.status=?", 3).
		Where("month=?", tools.ReturnTheMonth()).Scan(&result)

	//获取上月
	db.Table("task_orders").
		Select("sum(tasks.price) as last_month_get_money").Joins("LEFT JOIN tasks ON  task_orders.task_id=tasks.id").
		Where("task_orders.worker_id=?", t.WorkerId).Where("task_orders.status=?", 3).
		Where("month=?", tools.ReturnTheMonth()-1).Scan(&result)

	//获取总收益 all_get_money
	db.Table("task_orders").
		Select("sum(tasks.price) as all_get_money").Joins("LEFT JOIN tasks ON  task_orders.task_id=tasks.id").
		Where("task_orders.worker_id=?", t.WorkerId).Where("task_orders.status=?", 3).
		Scan(&result)

	//今日完成的订单数
	db.Model(TaskOrder{}).Where("status=?", 3).Where("worker_id=?", t.WorkerId).Where("date=?", time.Now().Format("2006-01-02")).Count(&result.TodayDoneTaskNum)
	//今日还剩下多少任务
	//获取用户 的当前 vip 等级

	result.TodayResidueTaskNum = taskTimes - result.TodayDoneTaskNum

	return result
}
