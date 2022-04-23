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
	"github.com/fatih/structs"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/wangyi/GinTemplate/tools"
	"time"
)

type AdminModel struct {
	ID         uint `gorm:"primaryKey;comment:'主键'"`
	Username   string
	Password   string
	Level      int
	RoleId     int    //角色id
	Status     int    // 1正常 2封禁
	GoogleCode string //谷歌  用来验证
	Token      string //48字符串
	Created    int64  //注册时间
}

func CheckIsExistModelAdminModel(db *gorm.DB, redis *redis.Client) {
	if db.HasTable(&AdminModel{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&AdminModel{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&AdminModel{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
		//创建原始数据
		token := tools.RandStringRunes(48)
		admin := AdminModel{ID: 1, Level: 0, Username: "admin", Password: "admin", Status: 1, Created: time.Now().Unix(), Token: token, RoleId: 1}
		err = db.Save(&admin).Error
		if err == nil {
			redis.HSet("Admin_Token", token, admin.Username).Result()
			redis.HMSet("Admin_"+admin.Username, structs.Map(&admin))
		}
	}
}
