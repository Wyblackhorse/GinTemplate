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

type RoleAndMenu struct {
	ID      uint  `gorm:"primaryKey;comment:'主键'"`
	RoleId  int   // 1 角色id
	MenuId  int   //菜单id
	Created int64 //创建时间
}

//初始化 超级管理员
func (r *RoleAndMenu) SuperAdminCreate(db *gorm.DB) {
	r.RoleId = 1
	r.Created = time.Now().Unix()
	err := db.Where("role_id=?", r.RoleId).Where("menu_id=?", r.MenuId).First(&RoleAndMenu{}).Error
	if err != nil {
		db.Save(&r)
	}
}

func CheckIsExistModelRoleAndMenu(db *gorm.DB) {
	if db.HasTable(&RoleAndMenu{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&RoleAndMenu{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&RoleAndMenu{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
			//初始化 数据
		}
	}
}
