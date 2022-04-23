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

//工人
type Worker struct {
	ID             uint    `gorm:"primaryKey;comment:'主键'"`
	Username       string  //用户名  真实姓名
	Password       string  //密码
	Balance        float64 `gorm:"type:decimal(10,2)"` //账户余额
	Token          string  //token
	InvitationCode string  //邀请码
	SuperiorId     int     `gorm:"int(10);default:0"` //上级id
	PayPassword    string  //资金密码
	HeadImage      string  //头像地址(给一个默认值)
	BankCardId     int     //银行卡 id
	AllIncome      float64 `gorm:"type:decimal(10,2)"`
	Phone          string  //手机号
	EMail          string  //电子邮箱
	Status         int     //状态 1限制  2良好 3优秀  4封号
	CreditScore    int     //信用积分
	Created        int64
}

func CheckIsExistModelWorker(db *gorm.DB) {
	if db.HasTable(&Role{}) {
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
