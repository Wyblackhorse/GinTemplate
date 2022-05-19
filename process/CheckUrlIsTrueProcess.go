package process

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/wangyi/GinTemplate/model"
	"net/http"
	"strconv"
	"time"
)

var sum int = 0

// SetCheckWebNameIsTrueProcess 生成  待检查的域名
func SetCheckWebNameIsTrueProcess(hash string, web model.WebName, redis *redis.Client, db *gorm.DB) {
	str := "0123456789abcdefghijklmnopqrstuvwxyz" //限定在8位
	//str="abc"
	for i := 1; i < 9-len(hash); i++ {
		Recursive("", i, str, hash, web, db)
	}
	//fmt.Printf("总数 %d", sum)

	//检查完毕
	redis.HDel("DOING", strconv.Itoa(int(web.ID)))

}

func Recursive(prefix string, total int, str string, webName string, web model.WebName, db *gorm.DB) {
	for i := 0; i < len(str); i++ {
		str1 := str //复制一份
		//str2 := str[0:i] + str[i+1:]
		str1 = str1[i : i+1]
		s := prefix + str1
		if total == 1 {
			//fmt.Println(s)
			InsertWebName(webName, s, web, db)
			sum++
		} else {
			temp := total - 1
			Recursive(s, temp, str, webName, web, db)
		}
	}
}

func InsertWebName(webName string, s string, web model.WebName, db *gorm.DB) {
	for i := 0; i < len(s)+1; i++ {
		//fmt.Println(webName+s[])
		var url string
		if i == 0 {
			url = webName + s
		} else if i == len(s) {
			url = s + webName
		} else {
			url = s[0:i] + webName + s[i:]
		}
		c := model.CheckWebName{WebNameId: int(web.ID), MatchUrl: url + "." + web.Suffix}
		c.AddCheckWebName(db)
	}
}

// CheckWebNameIsTrueProcess 检查 url 是否有效
func CheckWebNameIsTrueProcess(db *gorm.DB) {
	for true {
		c := make([]model.CheckWebName, 0)
		err2 := db.Where("status=?", 3).Limit(10000).Find(&c).Error
		if err2 == nil {
			for _, v := range c {
				_, err := http.Get("http://www." + v.MatchUrl)
				if err != nil {
					db.Model(&model.CheckWebName{}).Where("id=?", v.ID).Update(&model.CheckWebName{Updated: time.Now().Unix(), Status: 1}) //无效
					return
				}
				db.Model(&model.CheckWebName{}).Where("id=?", v.ID).Update(&model.CheckWebName{Updated: time.Now().Unix(), Status: 2}) //有效
				time.Sleep(1 * time.Second)
			}
		}

		time.Sleep(1 * time.Second)
	}

}
