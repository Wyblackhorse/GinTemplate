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

//工人
type Worker struct {
	ID                 uint    `gorm:"primaryKey;comment:'主键'"`
	Username           string  //用户名  真实姓名
	Password           string  //密码
	Balance            float64 `gorm:"type:decimal(10,2)"` //账户余额
	WithdrawalToFreeze float64 `gorm:"type:decimal(10,2)"` //提现冻结金额
	Token              string  //token
	InvitationCode     string  //邀请码
	SuperiorId         int     `gorm:"int(10);default:0"` //上级id
	PayPassword        string  //资金密码
	HeadImage          string  //头像地址(给一个默认值)
	BankCardId         int     //银行卡 id
	AllIncome          float64 `gorm:"type:decimal(10,2)"`
	YuEBaoMoney        float64 `gorm:"type:decimal(10,2)"` //余额宝里面的钱
	Phone              string  //手机号
	EMail              string  //电子邮箱
	Status             int     //状态 1限制  2良好 3优秀  4封号
	CreditScore        int     //信用积分
	VipId              int     `gorm:"int(10);default:1"` //vip  等级 id
	Created            int64
	VipName            string `gorm:"-"`
}

type WorkerBalance struct {
	ID              int     //用户的 id
	AddBalance      float64 //增加多少钱 减少
	ChangeMoneyLock sync.RWMutex
	Kinds           int //类型 1充值  2提现 3做单任务 4购买业务 5佣金奖励 6充值到余额宝
	OrderId         int
	YuEBaoId        int //余额宝产品id
	Days            int //理财产品的时间

}

func CheckIsExistModelWorker(db *gorm.DB) {
	if db.HasTable(&Worker{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Worker{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Worker{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}

//给用户加钱/扣金额
func (w *WorkerBalance) AddBalanceFuc(db *gorm.DB) (bool, error) {
	db = db.Begin() //开启事务
	//读锁
	w.ChangeMoneyLock.RLock()
	worker := Worker{}
	err := db.Where("id=?", w.ID).First(&worker).Error
	if err != nil {
		w.ChangeMoneyLock.RUnlock() //解除读锁
		return false, err
	}
	//读取正常  解除读锁
	w.ChangeMoneyLock.RUnlock()      //解除读锁
	w.ChangeMoneyLock.Lock()         //上写锁
	defer w.ChangeMoneyLock.Unlock() //解锁

	//加钱操作
	var newBalance float64
	if w.Kinds == 3 {
		newBalance = worker.Balance + w.AddBalance
		err = db.Model(&Worker{}).Where("id=?", w.ID).Update(&Worker{Balance: newBalance}).Error
		if err != nil {
			db.Rollback() //事务回滚
			return false, err
		}
	} else if w.Kinds == 6 || w.Kinds == 2 {
		//减钱操作
		if worker.Balance < w.AddBalance {
			//余额不够了
			return false, eeor.OtherError("don't have enough money")
		}
		newBalance := worker.Balance - w.AddBalance
		ps := map[string]interface{}{}
		ps["Balance"] = newBalance
		err = db.Model(&Worker{}).Where("id=?", w.ID).Update(ps).Error
		if err != nil {
			db.Rollback() //事务回滚
			return false, err
		}
	}
	//金额改变成功 类型 1充值  2提现 3做单任务 4购买业务 5佣金奖励
	if w.Kinds == 3 {
		add := BillingDetails{
			WorkerId:    int(worker.ID),
			ChangeMoney: w.AddBalance,
			InitMoney:   worker.Balance,
			NowMoney:    newBalance,
			Created:     time.Now().Unix(),
			Kinds:       w.Kinds,
		}
		err = db.Save(&add).Error
		if err != nil {
			db.Rollback() //事务回滚
			return false, err
		}
		//对订单进行更新
		err = db.Model(&TaskOrder{}).Where("id=?", w.OrderId).Update(&TaskOrder{Status: 3, Updated: time.Now().Unix()}).Error
		if err != nil {
			db.Rollback() //事务回滚
			return false, err
		}
	} else if w.Kinds == 6 {
		//充值到 余额宝   生成理财产品订单
		endTime := time.Now().Unix() + int64(w.Days*3600*24)
		addNewData := MoneyManagement{WorkerId: w.ID, YuEBaoId: w.YuEBaoId, Money: w.AddBalance, Status: 1, Created: time.Now().Unix(), EndTime: endTime}
		err = db.Save(&addNewData).Error
		if err != nil {
			db.Rollback() //事务回滚
			return false, err
		}
	} else if w.Kinds == 2 {
		//提现功能   增加冻结金额
		newWithdrawalToFreeze := worker.WithdrawalToFreeze + w.AddBalance
		err = db.Model(&Worker{}).Where("id=?", w.ID).Update(&Worker{WithdrawalToFreeze: newWithdrawalToFreeze}).Error
		if err != nil {
			db.Rollback() //事务回滚
			return false, err
		}
		//生成提现订单
		addRecord := Record{Kinds: 2, Status: 2, Money: w.AddBalance, WorkerId: int(worker.ID)}
		_, err = addRecord.AddRecord(db)
		if err != nil {
			db.Rollback() //事务回滚
			return false, err
		}
	}

	db.Commit() //事务提交
	return true, nil
}

//用户 增信用分
func (w *Worker) AddCreditScore(db *gorm.DB) {
	worker := Worker{}
	err := db.Where("id=?", w.ID).First(&worker).Error
	if err != nil {
		return
	}
	newCreditScore := worker.CreditScore + 1
	db.Model(&Worker{}).Where("id=?", w.ID).Update(map[string]interface{}{"CreditScore": newCreditScore})
}

//减信用分
func (w *Worker) SubtractCreditScore(db *gorm.DB) {
	worker := Worker{}
	err := db.Where("id=?", w.ID).First(&worker).Error
	if err != nil {
		return
	}
	newCreditScore := worker.CreditScore - 1
	db.Model(&Worker{}).Where("id=?", w.ID).Update(map[string]interface{}{"CreditScore": newCreditScore})
}

//判断用户是否存在
func (w *Worker) IsExist(db *gorm.DB) bool {
	err := db.Where("id=?", w.ID).First(&Worker{}).Error
	if err != nil {
		//
		return false
	}
	return true
}
