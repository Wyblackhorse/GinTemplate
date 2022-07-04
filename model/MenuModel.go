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
	IfChoose    int    `gorm:"-"` //是否选择  1勾选 2没有勾选
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
			m := Menu{ID: 1, Name: "首页", EnglishName: "homePage"}
			m.AddStairMenu(db)
			r := RoleAndMenu{ID: 1, MenuId: 1}
			r.SuperAdminCreate(db)

			m = Menu{ID: 2, Name: "会员管理", EnglishName: "memberManagement"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 2, MenuId: 2}
			r.SuperAdminCreate(db)

			//二级菜单  会员管理  2
			m = Menu{ID: 211, Name: "会员等级", Level: 2, EnglishName: "gradeOfMembership"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 211, MenuId: 211}
			r.SuperAdminCreate(db)
			m = Menu{ID: 212, Name: "普通会员", Level: 2, EnglishName: "regularMembers"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 212, MenuId: 212}
			r.SuperAdminCreate(db)

			m = Menu{ID: 213, Name: "会员银行", Level: 2, EnglishName: "memberBank"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 213, MenuId: 213}
			r.SuperAdminCreate(db)

			m = Menu{ID: 3, Name: "任务管理", EnglishName: "taskManagement"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 3, MenuId: 3}
			r.SuperAdminCreate(db)

			//二级菜单  任务管理 3
			m = Menu{ID: 311, Name: "任务列表", Level: 3, EnglishName: "taskList"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 311, MenuId: 311}
			r.SuperAdminCreate(db)

			m = Menu{ID: 312, Name: "任务种类", Level: 3, EnglishName: "taskKinds"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 312, MenuId: 312}
			r.SuperAdminCreate(db)

			m = Menu{ID: 313, Name: "审核任务", Level: 3, EnglishName: "reviewTheTask"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 313, MenuId: 313}
			r.SuperAdminCreate(db)

			m = Menu{ID: 314, Name: "采集任务", Level: 3, EnglishName: "acquisitionTask"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 314, MenuId: 314}
			r.SuperAdminCreate(db)

			m = Menu{ID: 4, Name: "账单管理", EnglishName: "billManagement"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 4, MenuId: 4}
			r.SuperAdminCreate(db)

			//二级菜单
			m = Menu{ID: 411, Name: "充值账单", Level: 4, EnglishName: "prepaidPhoneBills"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 411, MenuId: 411}
			r.SuperAdminCreate(db)

			m = Menu{ID: 412, Name: "提现账单", Level: 4, EnglishName: "withdrawalBill"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 412, MenuId: 412}
			r.SuperAdminCreate(db)

			m = Menu{ID: 413, Name: "佣金账单", Level: 4, EnglishName: "commissionBill"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 413, MenuId: 413}
			r.SuperAdminCreate(db)

			m = Menu{ID: 414, Name: "推广奖励", Level: 4, EnglishName: "promotionRewards"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 414, MenuId: 414}
			r.SuperAdminCreate(db)


			m = Menu{ID: 415, Name: "购买账单", Level: 4, EnglishName: "promotionRewards"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 415, MenuId: 415}
			r.SuperAdminCreate(db)

			m = Menu{ID: 5, Name: "报表管理", EnglishName: "statementManagement"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 5, MenuId: 5}
			r.SuperAdminCreate(db)

			//二级菜单
			m = Menu{ID: 511, Name: "每日报表", Level: 5, EnglishName: "statementEveryday"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 511, MenuId: 511}
			r.SuperAdminCreate(db)

			m = Menu{ID: 512, Name: "团队报表", Level: 5, EnglishName: "statementTeam"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 512, MenuId: 512}
			r.SuperAdminCreate(db)

			m = Menu{ID: 513, Name: "全局统计", Level: 5, EnglishName: "globalStatistics"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 513, MenuId: 513}
			r.SuperAdminCreate(db)

			m = Menu{ID: 6, Name: "财政管理", EnglishName: "fiscalManagement"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 6, MenuId: 6}
			r.SuperAdminCreate(db)

			m = Menu{ID: 7, Name: "应用管理", EnglishName: "applyManagement"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 7, MenuId: 7}
			r.SuperAdminCreate(db)

			m = Menu{ID: 711, Name: "应用列表", Level: 7, EnglishName: "applyList"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 711, MenuId: 711}
			r.SuperAdminCreate(db)

			m = Menu{ID: 8, Name: "团队管理", EnglishName: "teamManagement"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 8, MenuId: 8}
			r.SuperAdminCreate(db)

			m = Menu{ID: 811, Name: "团队列表", Level: 8, EnglishName: "teamList"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 811, MenuId: 811}
			r.SuperAdminCreate(db)

			m = Menu{ID: 9, Name: "日志管理", EnglishName: "logManagement"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 9, MenuId: 9}
			r.SuperAdminCreate(db)
			m = Menu{ID: 911, Name: "登录日志", Level: 9, EnglishName: "loginLog"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 911, MenuId: 911}
			r.SuperAdminCreate(db)
			m = Menu{ID: 912, Name: "系统日志", Level: 9, EnglishName: "systemLog"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 912, MenuId: 912}
			r.SuperAdminCreate(db)
			m = Menu{ID: 913, Name: "管理操作日志", Level: 9, EnglishName: "adminLog"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 913, MenuId: 913}
			r.SuperAdminCreate(db)

			m = Menu{ID: 10, Name: "权限管理", EnglishName: "jurisdictionManagement"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 10, MenuId: 10}
			r.SuperAdminCreate(db)

			m = Menu{ID: 1011, Name: "角色管理", Level: 10, EnglishName: "roleManagement"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 1011, MenuId: 1011}
			r.SuperAdminCreate(db)

			m = Menu{ID: 1012, Name: "权限分配", Level: 10, EnglishName: "roleAllocation"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 1012, MenuId: 1012}
			r.SuperAdminCreate(db)

			m = Menu{ID: 11, Name: "系统管理", EnglishName: "settingManagement"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 11, MenuId: 11}
			r.SuperAdminCreate(db)

			m = Menu{ID: 1111, Name: "基本设置", Level: 11, EnglishName: "basicSetting"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 1111, MenuId: 1111}
			r.SuperAdminCreate(db)

			m = Menu{ID: 1112, Name: "幻灯片设置", Level: 11, EnglishName: "lanternSlide"}
			m.AddSecondMenu(db)
			r = RoleAndMenu{ID: 1112, MenuId: 1112}
			r.SuperAdminCreate(db)

			m = Menu{ID: 12, Name: "余额宝", EnglishName: "balanceManagement"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 12, MenuId: 12}
			r.SuperAdminCreate(db)

			m = Menu{ID: 1211, Name: "产品列表", EnglishName: "productList"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 1211, MenuId: 1211}
			r.SuperAdminCreate(db)

			m = Menu{ID: 1212, Name: "购买记录", EnglishName: "purchaseHistory"}
			m.AddStairMenu(db)
			r = RoleAndMenu{ID: 1212, MenuId: 1212}
			r.SuperAdminCreate(db)

		}
	}
}
