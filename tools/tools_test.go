package tools

import (
	"fmt"
	"testing"
	"time"
)

func TestGetRootPath(t *testing.T) {
	////2022-07-01 13:00:00
	//loc, _ := time.LoadLocation("Asia/Shanghai")                   //设置时区
	//tt, _ := time.ParseInLocation("2006-01-02", "2022-07-01", loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
	//fmt.Println(tt.Unix())
	//tm := time.Unix(tt.Unix(), 0)
	//fmt.Println(tm.Format("2006-01-02 15:04:05"))
	//str := "2022-07-01"
	//fmt.Println(strings.Split(str, " ")[0])

	loc, _ := time.LoadLocation("Asia/Shanghai")                        //设置时区
	tt, _ := time.ParseInLocation("2006-01-02", "2022-07-04", loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"

	fmt.Println(tt.Unix())

	type oNE struct {

		Name  string


	}








}
