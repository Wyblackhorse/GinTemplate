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
	"math/rand"
	"strconv"
	"time"
)



type Record struct {
	ID               uint    `gorm:"primaryKey;comment:'主键'"`
	WorkerId         int     //用户id
	Kinds            int     //类型 1充值  2提现   4购买业务 5佣金奖励(邀请奖励) 6充值到余额宝  7任务提成(团队)
	Money            float64 `gorm:"type:decimal(10,2);default:0"` //购买金额
	AccountPractical float64 `gorm:"type:decimal(10,2);default:-1"` //充值订单的时候  实际充值的金额   -1 管理员操作 -2代表用户操作
	Status           int    //1已完成  2审核中  3失败
	Month            int    //月
	Week             int    //周
	Date             string //日期
	RecordNum        string
	Created          int64
	Updated          int64
	WorkerName       string `gorm:"-"`
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
	r.Date = time.Now().Format("2006-01-02")
	r.Month = tools.ReturnTheMonth()
	r.Week = tools.ReturnTheWeek()
	r.AccountPractical =-2
	r.RecordNum = time.Now().Format("20060102") + strconv.FormatFloat(float64(time.Now().Unix()), 'f', 0, 64) + strconv.Itoa(r.WorkerId)
	err := db.Save(&r).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

//是否存在这个订单  如果存在就返回 订单数据
func (r *Record) IsExistRecord(db *gorm.DB) (bool, Record) {
	err := db.Where("id=?", r.ID).First(&r).Error
	if err != nil {
		return false, Record{}
	}
	return true, *r
}

//管理员审核订单
func (r *Record) WithdrawDeposit(db *gorm.DB, status int) (bool, error) {
	//通过
	if status == 1 {
		//修改订单状态
		err := db.Model(&Record{}).Where("id=?", r.ID).Update(&Record{Status: 1, Updated: time.Now().Unix()}).Error
		if err != nil {
			//更新失败
			return false, err
		}
		worker := WorkerBalance{ID: r.WorkerId, AddBalance: r.Money, Kinds: 8}
		resultBool, err := worker.AddBalanceFuc(db)
		if resultBool == false {
			return false, err
		}

	} else if status == 3 {
		//不通过
		err := db.Model(&Record{}).Where("id=?", r.ID).Update(&Record{Status: 3, Updated: time.Now().Unix()}).Error
		if err != nil {
			//更新失败
			return false, err
		}

		worker := WorkerBalance{ID: r.WorkerId, AddBalance: r.Money, Kinds: 7}
		resultBool, err := worker.AddBalanceFuc(db)
		if resultBool == false {
			return false, err
		}
	}
	return true, nil
}

//创建充值订单
func (r *Record) CratedNewRechargeOrder(db *gorm.DB) string {
	for i := 0; i < 5; i++ {
		//首先生成一个订单号
		r.RecordNum = "DzCz" + time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(1000000))
		err := db.Where("record_num=?", r.RecordNum).First(&Record{}).Error
		if err != nil {
			r.Status = 2
			r.Month = tools.ReturnTheMonth()
			r.Week = tools.ReturnTheWeek()
			r.Date = time.Now().Format("2006-01-02")
			r.Updated = time.Now().Unix()
			r.Created = time.Now().Unix()
			err = db.Save(&r).Error
			if err != nil {
				continue
			}
			return r.RecordNum
		}
	}
	return ""
}
