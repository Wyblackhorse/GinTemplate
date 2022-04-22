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

type Loggers struct {
	ID      uint   `gorm:"primaryKey;comment:'主键'"`
	WriteId int    //操作者  如果为0  就是系统日志
	Kinds   int    //1玩家 2管理员
	Content string `gorm:"type:text"` //日志内容
	Status  int    //日志状态  1正常日志 2错误日志
	Created int64
}



func CheckIsExistModelLoggers(db *gorm.DB ) {
	if db.HasTable(&AdminModel{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Loggers{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Loggers{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}


//插入正常日志
func (log *Loggers) AdminAddNormalLoggers(Db *gorm.DB) {
	log.Status = 1
	log.Kinds = 2
	log.Created = time.Now().Unix()
	Db.Save(&log)
}

//插入错误日志
func (log *Loggers) AdminAddErrorLoggers(Db *gorm.DB) {
	log.Status = 2
	log.Kinds = 2
	log.Created = time.Now().Unix()
	Db.Save(&log)
}
