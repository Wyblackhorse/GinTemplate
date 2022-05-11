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
	NextSuperiorId     int     `gorm:"int(10);default:0"` //次上级id
	NextNextSuperiorId int     `gorm:"int(10);default:0"` //次上上级id
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
	Kinds           int //类型 1充值  2用户提现 3做单任务 4购买业务 5佣金奖励(邀请奖励) 6充值到余额宝   7管理员审核提现失败   8管理员审核提现成功
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

//给用户加钱/扣金额(余额变化表)
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
	if w.Kinds == 3 || w.Kinds == 7 || w.Kinds == 5 {
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

	//金额改变成功 类型 1充值  2提现 3做单任务 4购买业务 5佣金奖励(邀请奖励)
	if w.Kinds == 3 || w.Kinds == 5 {
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
		if w.Kinds == 3 {
			//对订单进行更新
			err = db.Model(&TaskOrder{}).Where("id=?", w.OrderId).Update(&TaskOrder{Status: 3, Updated: time.Now().Unix()}).Error
			if err != nil {
				db.Rollback() //事务回滚
				return false, err
			}
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
	} else if w.Kinds == 8 || w.Kinds == 7 {
		//管理员审核提现订单成功  余额不动  扣除 冻结提现金额
		newBalance := worker.WithdrawalToFreeze - w.AddBalance
		ps := map[string]interface{}{}
		ps["WithdrawalToFreeze"] = newBalance

		//审核失败  要把钱返回给余额
		if w.Kinds == 7 {
			newBalanceTwo := worker.Balance + w.AddBalance
			ps["Balance"] = newBalanceTwo
			add := BillingDetails{
				WorkerId:    int(worker.ID),
				ChangeMoney: w.AddBalance,
				InitMoney:   worker.Balance,
				NowMoney:    newBalanceTwo,
				Created:     time.Now().Unix(),
				Kinds:       w.Kinds,
			}
			err = db.Save(&add).Error
			if err != nil {
				db.Rollback() //事务回滚
				return false, err
			}

		}
		err = db.Model(&Worker{}).Where("id=?", w.ID).Update(ps).Error
		if err != nil {
			db.Rollback() //事务回滚
			return false, err
		}

	}
	db.Commit() //事务提交
	return true, nil
}

//用户 增信用分
func (r *Worker) AddCreditScore(db *gorm.DB) {
	worker := Worker{}
	err := db.Where("id=?", r.ID).First(&worker).Error
	if err != nil {
		return
	}
	newCreditScore := worker.CreditScore + 1
	db.Model(&Worker{}).Where("id=?", r.ID).Update(map[string]interface{}{"CreditScore": newCreditScore})
}

//减信用分
func (r *Worker) SubtractCreditScore(db *gorm.DB) {
	worker := Worker{}
	err := db.Where("id=?", r.ID).First(&worker).Error
	if err != nil {
		return
	}
	newCreditScore := worker.CreditScore - 1
	db.Model(&Worker{}).Where("id=?", r.ID).Update(map[string]interface{}{"CreditScore": newCreditScore})
}

//判断用户是否存在
func (r *Worker) IsExist(db *gorm.DB) bool {
	err := db.Where("id=?", r.ID).First(&Worker{}).Error
	if err != nil {
		//
		return false
	}
	return true
}

//获取团队信息

type TeamInformation struct {
	TeamRecharge        float64 `json:"team_recharge"`         //团队充值
	TeamWithdraw        float64 `json:"team_withdraw"`         //团队提现
	FirstRechargePeople int     `json:"first_people"`          //首充人数
	FirstPushPeople     int     `json:"first_push_people"`     //首推人数
	TeamPeople          int     `json:"team_people"`           //团队人数
	TeamNewAdd          int     `json:"team_new_add"`          //团队新增
	OneRecharge         float64 `json:"one_recharge"`          //一级充值金额
	OneRechargePeople   int     `json:"one_recharge_people"`   //一级充值人数
	OneBackMoney        float64 `json:"one_back_money"`        //一级返佣金额
	TwoRecharge         float64 `json:"two_recharge"`          //二级充值金额
	TwoRechargePeople   int     `json:"two_recharge_people"`   //二级充值人数
	TwoBackMoney        float64 `json:"two_back_money"`        //二级返佣金额
	ThreeRecharge       float64 `json:"three_recharge"`        //三级充值金额
	ThreeRechargePeople int     `json:"three_recharge_people"` //三级充值人数
	ThreeBackMoney      float64 `json:"three_back_money"`      //三级返佣金额
}

//团队报表
func (r *Worker) GetTeamStatement(db *gorm.DB) TeamInformation {

	var returnJson TeamInformation
	//一级 (充值的金额)
	db.Table("workers").
		Select("sum(records.money) as one_recharge").Joins("LEFT JOIN records ON  records.worker_id=workers.id").
		Where("workers.superior_id=?", r.ID).Where(" records.status=?", 1).Scan(&returnJson)
	//一级(充值人数)
	db.Raw("select count(*) as one_recharge_people  from( SELECT  1  FROM workers     LEFT JOIN records on  records.worker_id=workers.id WHERE records.status=1  AND workers.superior_id=? AND records.kinds=1    GROUP BY workers.superior_id )  t", r.ID).Scan(&returnJson)
	//一级(返佣金额)
	returnJson.OneBackMoney = 0
	//二级
	db.Table("workers").
		Select("sum(records.money) as two_recharge").Joins("LEFT JOIN records ON  records.worker_id=workers.id").
		Where("workers.next_superior_id=?", r.ID).Where(" records.status=?", 1).Scan(&returnJson)
	db.Raw("select count(*) as two_recharge_people  from( SELECT  1  FROM workers     LEFT JOIN records on  records.worker_id=workers.id WHERE records.status=1  AND workers.superior_id=? AND records.kinds=1    GROUP BY workers.next_superior_id )  t", r.ID).Scan(&returnJson)
	returnJson.TwoBackMoney = 0
	//三级
	db.Table("workers").
		Select("sum(records.money) as three_recharge").Joins("LEFT JOIN records ON  records.worker_id=workers.id").
		Where("workers.next_next_superior_id=?", r.ID).Where(" records.status=?", 1).Scan(&returnJson)
	db.Raw("select count(*) as three_recharge_people  from( SELECT  1  FROM workers     LEFT JOIN records on  records.worker_id=workers.id WHERE records.status=1  AND workers.superior_id=? AND records.kinds=1    GROUP BY workers.next_next_superior_id )  t", r.ID).Scan(&returnJson)
	returnJson.ThreeBackMoney = 0
	//团队充值
	returnJson.TeamRecharge = returnJson.OneRecharge + returnJson.TwoBackMoney + returnJson.ThreeBackMoney
	//团队提现
	db.Raw("SELECT  SUM(records.money) as team_withdraw  FROM workers  LEFT JOIN records ON records.worker_id=workers.id  WHERE   workers.superior_id= ? or  workers.next_superior_id= ?  or  workers.next_next_superior_id=?   AND  records.kinds=2", r.ID, r.ID, r.ID).Scan(&returnJson)
	//首充人数
	returnJson.FirstRechargePeople = returnJson.TwoRechargePeople + returnJson.OneRechargePeople + returnJson.TwoRechargePeople
	//团队的 人数
	db.Raw("SELECT  count(*) as team_people  FROM workers WHERE  superior_id=?  or  next_superior_id=? or  next_next_superior_id=?", r.ID, r.ID, r.ID).Scan(&returnJson)
	//1级首推
	db.Raw("SELECT COUNT(*) as first_push_people   from  (SELECT 1  FROM  workers  WHERE  next_superior_id=?    GROUP BY superior_id)  t", r.ID).Scan(&returnJson)
	one := returnJson.FirstPushPeople
	db.Raw("SELECT COUNT(*) as first_push_people   from  (SELECT 1  FROM  workers  WHERE  next_next_superior_id=?    GROUP BY next_superior_id)  t", r.ID).Scan(&returnJson)
	returnJson.FirstPushPeople = one + returnJson.FirstPushPeople

	return returnJson

}

type MyTeamInformation struct {
	Name      string  `json:"name"`
	Recharge  float64 `json:"recharge"`   //充值
	Withdraw  float64 `json:"withdraw"`   //提现
	BackMoney float64 `json:"back_money"` //反点
	Royalties float64 `json:"royalties"`  //提成
}

type ReturnMyTeamInformation struct {
	data []MyTeamInformation
}

//我的团队
func (r *Worker) GetMyTeamInformation(db *gorm.DB) []MyTeamInformation {

	//用户  充值   提现  返点 提成
	w := make([]Worker, 0)
	var data ReturnMyTeamInformation
	//获取下级
	db.Where("superior_id=? or  next_superior_id=?  or next_next_superior_id=?", r.ID, r.ID, r.ID).Find(&w)
	for _, v := range w {
		var a MyTeamInformation
		a.Name = v.Phone
		//充值
		db.Raw("select SUM(money)  as  recharge from records where kinds=1 AND worker_id=?", v.ID).Scan(&a)
		//提现
		db.Raw("select SUM(money)  as  recharge from records where kinds=2 AND worker_id=?", v.ID).Scan(&a)
		//返点
		a.BackMoney = 0
		//提成
		a.Royalties = 0
		data.data = append(data.data, a)
	}
	return data.data
}
