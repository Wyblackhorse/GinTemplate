package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	eeor "github.com/wangyi/GinTemplate/error"
	"github.com/wangyi/GinTemplate/util"
	"math/rand"
	"strconv"
	"time"
)

// PrepaidPhoneOrders 充值订单
type PrepaidPhoneOrders struct {
	ID                uint    `gorm:"primaryKey;comment:'主键'"`
	PlatformOrder     string  //平台订单  (前台订单号)
	ThreeOrder        string  //三方订单  (自己 随机成功)
	RechargeAddress   string  //充值地址
	CollectionAddress string  //收款地址
	RechargeType      string  //充值类型
	Username          string  //充值用户名
	AccountOrders     float64 `gorm:"type:decimal(10,2)"` //充值金额 (订单金额)
	AccountPractical  float64 `gorm:"type:decimal(10,2)"` //充值金额(实际返回金额)
	Status            int     //订单状态  1 未支付  2已经支付了  3已经失效
	ThreeBack         int     //三方回调 1未回调  2已结回调
	Created           int64   //订单创建时间
	Updated           int64   //更新时间(回调时间)
	Successfully      int64   //交易成功 时间(区块时间戳)
	Date              string
	BackUrl           string //回调的地址
}

func CheckIsExistModePrepaidPhoneOrders(db *gorm.DB) {
	if db.HasTable(&PrepaidPhoneOrders{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&PrepaidPhoneOrders{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&PrepaidPhoneOrders{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}

// CreatePrepaidPhoneOrders 创建充值订单
func (p *PrepaidPhoneOrders) CreatePrepaidPhoneOrders(db *gorm.DB) (bool, error) {
	//看是否存在有效的订单
	pt := make([]PrepaidPhoneOrders, 0)
	err2 := db.Where("username=?", p.Username).Where("status=?", 1).Find(&pt).Error
	if err2 == nil {
		//修改状态
		for _, v := range pt {
			db.Model(&PrepaidPhoneOrders{}).Where("id=?", v.ID).Update(&PrepaidPhoneOrders{Status: 3})
		}
	}
	p.Created = time.Now().Unix()
	p.Updated = 0
	if p.ThreeOrder == "" {
		p.ThreeOrder = time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(100000))
	}

	p.Status = 1
	p.ThreeBack = 1
	p.Date = time.Now().Format("2006-01-02")
	//创建之前判断是否有事重复提交
	err := db.Where("platform_order=?", p.PlatformOrder).First(&PrepaidPhoneOrders{}).Error
	if err == nil {
		return false, eeor.OtherError("重复提交")
	}
	err = db.Save(&p).Error
	if err != nil {
		return false, err
	}
	return true, nil

}

// UpdateMaxCreatedOfStatusToTwo 寻找一条最新传建的订单并且修改他的状态
func (p *PrepaidPhoneOrders) UpdateMaxCreatedOfStatusToTwo(db *gorm.DB, OrderEffectivityTime int64) bool {
	//找到这条数据
	pp := PrepaidPhoneOrders{}
	err := db.Where("username=?", p.Username).Where("status= ? and recharge_type= ?", 1,p.RechargeType).Last(&pp).Error
	if err == nil {
		if time.Now().Unix()-pp.Created <= OrderEffectivityTime {
			//找到最新的数据(并且在有效时间累)
			db.Model(&PrepaidPhoneOrders{}).Where("id=?", pp.ID).Update(
				&PrepaidPhoneOrders{Updated: time.Now().Unix(), Successfully: p.Successfully, ThreeBack: 2, Status: 2, AccountPractical: p.AccountPractical})

			//这里 要回调给前台
			if pp.BackUrl != "" {
				type Create struct {
					PlatformOrder    string
					RechargeAddress  string
					Username         string
					AccountOrders    float64 //订单充值金额
					AccountPractical float64 //  实际充值的金额
					RechargeType     string
					BackUrl          string
				}

				var tt Create
				tt.PlatformOrder = pp.PlatformOrder
				tt.RechargeAddress = p.RechargeAddress
				tt.Username = p.Username
				tt.AccountOrders = pp.AccountOrders
				tt.AccountPractical = p.AccountPractical
				tt.RechargeType = p.RechargeType


				data, err := json.Marshal(tt)
				if err != nil {
					return false
				}
				data, err= util.RsaEncryptForEveryOne(data)

				util.BackUrlToPay(pp.BackUrl, base64.StdEncoding.EncodeToString(data))
			}

			return true
		} else {
			db.Model(&PrepaidPhoneOrders{}).Where("id=?", pp.ID).Update(&PrepaidPhoneOrders{Status: 3})
			return false
		}

	}
	//没有这条数据  默认是线下支付   自行创建这笔订单 Username: p.UserID, Successfully: p.Timestamp, AccountPractical: p.Amount}
	pt := PrepaidPhoneOrders{}
	pt.Created = time.Now().Unix()
	pt.Updated = 0
	pt.ThreeOrder = time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(100000))
	pt.PlatformOrder = ""
	pt.Status = 2
	pt.ThreeBack = 3
	pt.Username = p.Username
	pt.Successfully = p.Successfully
	pt.AccountPractical = p.AccountPractical
	pt.RechargeType = p.RechargeType
	pt.RechargeAddress = p.RechargeAddress
	pt.CollectionAddress = p.CollectionAddress
	pt.Date = time.Now().Format("2006-01-02")
	db.Save(&pt)
	return false

}
