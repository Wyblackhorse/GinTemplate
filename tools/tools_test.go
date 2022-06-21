/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package tools

import (
	"fmt"
	"testing"
	"time"
)

func TestRandStringRunes(t *testing.T) {

	fmt.Println("-----------------开启二次认证----------------------")
	user := "testxxx@qq.com"
	secret, qrCodeUrl, a := InitAuth(user)

	fmt.Println("-----------------信息校验----------------------")

	fmt.Println(secret, qrCodeUrl, a)

	// secret最好持久化保存在
	//// 验证,动态码(从谷歌验证器获取或者freeotp获取)
	//bool, err := NewGoogleAuth().VerifyCode("QNSCBEQ4M6MAOIWB3H2NZYMN3JUD3SVT", "899397")
	//if bool {
	//	fmt.Println("√")
	//} else {
	//	fmt.Println("X", err)
	//}

}

func TestOne(t *testing.T) {

	fmt.Println(time.Now().AddDate(0, -1, 0).Format("2006-01-02"))

	fmt.Println(time.Now().Date())
	_, m, _ := time.Now().Date()
	fmt.Println(int(m))

	fmt.Println(time.Now().Weekday())
	datetime := time.Now().Format("20060102")
	timeLayout := "20060102"
	loc, _ := time.LoadLocation("Local")
	tmp, _ := time.ParseInLocation(timeLayout, datetime, loc)
	_, intWeek := tmp.ISOWeek()
	fmt.Println(intWeek)
}

func TestTwo(t *testing.T) {

	type ReturnData struct {
		TaskNum        int     `json:"task_num"`        //任务数量
		TaskEarnings   float64 `json:"task_earnings"`   //任务收益
		JuniorTaskNum  int     `json:"junior_task_num"` //下级任务数量
		JuniorEarnings float64 `json:"junior_earnings"` //下级任务收益
		Date           string  `json:"date"`
	}

	//c := make([]ReturnData, 0)
	//c1 := ReturnData{TaskNum: 10, TaskEarnings: 100, Date: "6-16"}
	//c2 := ReturnData{TaskNum: 9, TaskEarnings: 10, Date: "6-12"}
	//c3 := ReturnData{TaskNum: 8, TaskEarnings: 190, Date: "6-13"}
	//
	//
	//
	//b1 := ReturnData{JuniorTaskNum: 10, JuniorEarnings: 100, Date: "6-16"}
	//b2 := ReturnData{JuniorTaskNum: 9, JuniorEarnings: 10, Date: "6-12"}
	//b3 := ReturnData{JuniorTaskNum: 8, JuniorEarnings: 190, Date: "6-10"}
	//b4 := ReturnData{JuniorTaskNum: 8, JuniorEarnings: 200, Date: "6-17"}

}

func TestFormYouTuBeUrl(t *testing.T) {

	a := []int{1, 2, 4}

	var  b  [][]int

	b=append(b, a)



	b=append(b, a)
	fmt.Println(b)

}
