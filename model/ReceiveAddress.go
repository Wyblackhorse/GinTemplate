package model

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/wangyi/GinTemplate/tools"
	"go.uber.org/zap"
	"time"
)

// ReceiveAddress 收账地址管理
type ReceiveAddress struct {
	ID             uint `gorm:"primaryKey;comment:'主键'"`
	Username       string
	ReceiveNums    int     //收款笔数
	LastGetAccount float64 `gorm:"type:decimal(10,2)"` //最后一次的入账金额
	Address        string  //收账地址
	Money          float64 `gorm:"type:decimal(10,2)"` //账户余额
	Created        int64
	Updated        int64
}

func CheckIsExistModeReceiveAddress(db *gorm.DB) {
	if db.HasTable(&ReceiveAddress{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&ReceiveAddress{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&ReceiveAddress{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}

// ReceiveAddressIsExits 判断转账地址是否存在
func (r *ReceiveAddress) ReceiveAddressIsExits(db *gorm.DB) bool {
	err := db.Where("username=?", r.Username).First(&ReceiveAddress{}).Error
	if err != nil {
		//错误存在(没有这个用户)
		return false
	}
	return true
}

// CreateUsername 创建这个用户  获取用户收款地址
func (r *ReceiveAddress) CreateUsername(db *gorm.DB, url string)  ReceiveAddress{
	r.Created = time.Now().Unix()
	r.ReceiveNums = 0
	r.LastGetAccount = 0
	//获取收账地址  url 请求  {"error":"0","message":"","result":"4564554545454545"}   //返回数据
	req := make(map[string]interface{})
	req["user"] = r.Username
	req["ts"] = time.Now().UnixMilli()
	resp, err := tools.HttpRequest(url+"/getaddr", req, viper.GetString("eth.ApiKey"))
	if err != nil {
		fmt.Println(err.Error())
		return ReceiveAddress{}
	}
	var dataAttr CreateUsernameData
	if err := json.Unmarshal([]byte(resp), &dataAttr); err != nil {
		fmt.Println(err)
		return ReceiveAddress{}
	}
	if dataAttr.Result != "" {
		r.Address = dataAttr.Result
		err := db.Save(&r).Error
		if err!=nil {
			return  ReceiveAddress{}
		}
	}
	return *r
}

// CreateUsernameData 返回的数据 json
type CreateUsernameData struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func (r *ReceiveAddress) UpdateReceiveAddressLastInformation(db *gorm.DB) bool {
	re := ReceiveAddress{}
	err := db.Where("username=?", r.Username).First(&re).Error
	if err == nil {
		nums := re.ReceiveNums + 1
		err := db.Model(&ReceiveAddress{}).Where("id=?", re.ID).Update(&ReceiveAddress{ReceiveNums: nums, LastGetAccount: r.LastGetAccount, Updated: r.Updated, Money: r.Money}).Error
		if err == nil {
			return true
		}
	}
	return false
}

func (r *ReceiveAddress) UpdateReceiveAddressLastInformationTo0(db *gorm.DB) bool {
	re := ReceiveAddress{}
	err := db.Where("username=?", r.Username).First(&re).Error
	if err == nil {
		zap.L().Debug("余额清0,用户:" + r.Username )

		updated := make(map[string]interface{})
		updated["Updated"] = r.Updated
		updated["Money"] = 0
		err := db.Model(&ReceiveAddress{}).Where("id=?", re.ID).Update(updated).Error
		if err == nil {
			return true
		}
	}

	zap.L().Debug("余额清0,用户:" + r.Username + "没有找到")

	return false
}

// CreateNewReceiveAddress 创建新的地址
func CreateNewReceiveAddress(db *gorm.DB, url string) bool {
	//随机生成新的用户名
	username := tools.RandString(40)
	err := db.Where("username=?", string(username)).First(&ReceiveAddress{}).Error
	if err == nil {
		//找到了
		return false
	}
	r2 := ReceiveAddress{Username: string(username)}
	r2.CreateUsername(db, url)
	return true
}
