package model

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/wangyi/GinTemplate/tools"
	"io/ioutil"
	"net/http"
	"time"
)

// ReceiveAddress 收账地址管理
type ReceiveAddress struct {
	ID             uint `gorm:"primaryKey;comment:'主键'"`
	Username       string
	ReceiveNums    int    //收款笔数
	LastGetAccount int    //最后一次的入账金额
	Address        string //收账地址
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

// CreateUsername 创建这个用户
func (r *ReceiveAddress) CreateUsername(db *gorm.DB, url string) {
	r.Created = time.Now().Unix()
	r.Updated = time.Now().Unix()
	r.ReceiveNums = 0
	r.LastGetAccount = 0
	//获取收账地址  url 请求  {"error":"0","message":"","result":"4564554545454545"}   //返回数据
	fmt.Println(url + "/getaddr?user=" + r.Username)
	resp, err := http.Get(url + "/getaddr?user=" + r.Username)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	var dataAttr CreateUsernameData
	if err := json.Unmarshal([]byte(body), &dataAttr); err != nil {
		fmt.Println(err)
		return
	}
	if dataAttr.Result != "" {
		r.Address = dataAttr.Result
		db.Save(&r)
	}

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
		err := db.Model(&ReceiveAddress{}).Where("id=?", re.ID).Update(&ReceiveAddress{ReceiveNums: nums, LastGetAccount: r.LastGetAccount, Updated: time.Now().Unix()}).Error
		if err == nil {
			return true
		}

	}
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
