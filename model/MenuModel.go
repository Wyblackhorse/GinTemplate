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

type Menu struct {
	ID          uint   `gorm:"primaryKey;comment:'主键'"`
	Name        string //菜单名称
	EnglishName string //英文名字  用于权限控制
	Level       int    //菜单等级 0 一级菜单
	Status      int    //1正常 2禁用
	Created     int64  //创建时间
}

//添加一级菜单
func (m *Menu) AddStairMenu(Db *gorm.DB) {
	//判断是否重复添加
	m.Level = 0
	m.Status = 1
	m.Created = time.Now().Unix()
	err := Db.Where("name=?", m.Name).First(&Menu{}).Error
	if err != nil {
		Db.Save(&m)
	}
}

//添加二级菜单
func (m *Menu) AddSecondMenu(Db *gorm.DB) {
	//判断是否重复添加
	m.Status = 1
	m.Created = time.Now().Unix()
	err := Db.Where("name=?", m.Name).First(&Menu{}).Error
	if err != nil {
		Db.Save(&m)
	}
}

func CheckIsExistModelMenu(db *gorm.DB) {
	if db.HasTable(&Menu{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Menu{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Menu{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
			//初始化数据
			m := Menu{ID: 1, Name: "首页"}
			m.AddStairMenu(db)
			r := RoleAndMenu{ID: 1, MenuId: 1}
			r.SuperAdminCreate(db)

			m = Menu{ID: 2, Name: "会员管理"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 2, MenuId: 2}
			r.SuperAdminCreate(db)

			m = Menu{ID: 3, Name: "任务管理"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 3, MenuId: 3}
			r.SuperAdminCreate(db)

			m = Menu{ID: 4, Name: "账单管理"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 4, MenuId: 4}
			r.SuperAdminCreate(db)

			m = Menu{ID: 5, Name: "报表管理"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 5, MenuId: 5}
			r.SuperAdminCreate(db)

			m = Menu{ID: 6, Name: "财政管理"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 6, MenuId: 6}
			r.SuperAdminCreate(db)

			m = Menu{ID: 7, Name: "应用管理"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 7, MenuId: 7}
			r.SuperAdminCreate(db)

			m = Menu{ID: 8, Name: "团队管理"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 8, MenuId: 8}
			r.SuperAdminCreate(db)

			m = Menu{ID: 9, Name: "日志管理"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 9, MenuId: 9}
			r.SuperAdminCreate(db)

			m = Menu{ID: 10, Name: "权限管理"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 10, MenuId: 10}
			r.SuperAdminCreate(db)

			m = Menu{ID: 11, Name: "系统管理"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 11, MenuId: 11}
			r.SuperAdminCreate(db)

			//二级菜单  会员管理  2
			m = Menu{ID: 121, Name: "会员等级", Level: 2}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 121, MenuId: 121}
			r.SuperAdminCreate(db)
			m = Menu{ID: 122, Name: "普通会员", Level: 2}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 122, MenuId: 122}
			r.SuperAdminCreate(db)

			//任务管理 3
			m = Menu{ID: 131, Name: "已完成", Level: 3}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 131, MenuId: 131}
			r.SuperAdminCreate(db)




		}
	}
}
