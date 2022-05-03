package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	eeor "github.com/wangyi/GinTemplate/error"
	"math/rand"
	"strconv"
	"time"
)

// PrepaidPhoneOrders 充值订单
type PrepaidPhoneOrders struct {
	ID                uint   `gorm:"primaryKey;comment:'主键'"`
	PlatformOrder     string //平台订单
	ThreeOrder        string //三方订单
	RechargeAddress   string //充值地址
	CollectionAddress string //收款地址
	RechargeType      string //充值类型
	Username          string //充值用户名
	AccountOrders     int    //充值金额 (订单金额)
	AccountPractical  int    //充值金额(实际返回金额)
	Status            int    //订单状态  1 未支付  2已经支付了
	ThreeBack         int    //三方回调 1未回调  2已结回调
	Created           int64  //订单创建时间
	Updated           int64  //更新时间(回调时间)
	Successfully      int64  //交易成功 时间(区块时间戳)
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
	p.Created = time.Now().Unix()
	p.Updated = 0
	p.PlatformOrder = time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(100000))
	p.Status = 1
	p.ThreeBack = 1
	//创建之前判断是否有事重复提交
	err := db.Where("three_order=?", p.ThreeOrder).First(&PrepaidPhoneOrders{}).Error
	if err == nil {
		return false, eeor.OtherError("重复提交")
	}
	err = db.Save(&p).Error
	if err != nil {
		return false, err
	}
	return true, nil

}

//寻找一条最新传建的订单并且修改他的状态

func (p *PrepaidPhoneOrders) UpdateMaxCreatedOfStatusToTwo(db *gorm.DB) bool {
	//找到这条数据
	pp := PrepaidPhoneOrders{}
	err := db.Where("username=?", p.Username).Last(&pp).Error
	if err == nil {
		err := db.Model(&PrepaidPhoneOrders{}).Where("id=?", pp.ID).Update(
			&PrepaidPhoneOrders{Updated: time.Now().Unix(), Successfully: p.Successfully, ThreeBack: 2, Status: 2, AccountPractical: p.AccountPractical}).Error
		if err != nil {
			return false
		}
		return true

	}

	return false

}
