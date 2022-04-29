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
	"time"
)

type YuEBao struct {
	ID           uint    `gorm:"primaryKey;comment:'主键'"`
	Name         string  //产品名字
	MinMoney     float64 `gorm:"type:decimal(10,2)"` //最低购买
	InterestRate float64 `gorm:"type:decimal(10,2)"` //利率
	Days         int     // 时间
	Status       int     //状态    1开启 2关闭
	Created      int64   //创建时间

}

func CheckIsExistModelYuEBao(db *gorm.DB) {
	if db.HasTable(&YuEBao{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&YuEBao{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&YuEBao{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}

//添加数据
func (y *YuEBao) AddYuEBao(db *gorm.DB) (bool, error) {
	y.Created = time.Now().Unix()
	y.Status = 1
	err := db.Where("name=?", y.Name).First(&YuEBao{}).Error
	if err == nil {
		return false, nil
	}
	err = db.Save(&y).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
